package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"gorag/services"
)

type askRequestBody struct {
	Question string `json:"question"`
	Source   string `json:"source"`
}

type ollamaGenerate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func GetAskHandler(client *weaviate.Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(400)
			return
		}

		var r askRequestBody
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer req.Body.Close()

		embedding, err := services.GetEmbedding(r.Question)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		results, err := services.SimilaritySearch(client, embedding, r.Source, 5)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}

		context := strings.Join(results, "\n\n")
		body := ollamaGenerate{
			Model: "qwen2.5:7b",
			Prompt: fmt.Sprintf(`Answer the question only based on the following context:

Context: "%s"

Question: "%s"
`, context, r.Question),
		}
		data, err := json.Marshal(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}

		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}

		resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(data))
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	}
}
