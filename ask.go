package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

type AskRequestBody struct {
	Question string `json:"question"`
	Source   string `json:"source"`
}

type OllamaGenerate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func similaritySearch(client *weaviate.Client, queryVector []float32, filename string, limit int) ([]string, error) {
	nearVector := client.GraphQL().NearVectorArgBuilder().
		WithVector(queryVector)

	filter := filters.Where().
		WithPath([]string{"source"}).
		WithOperator(filters.Equal).
		WithValueText(filename)

	result, err := client.GraphQL().Get().
		WithClassName("Document").
		WithFields(
			graphql.Field{Name: "text"},
			graphql.Field{Name: "source"},
			graphql.Field{
				Name: "_additional",
				Fields: []graphql.Field{
					{Name: "distance"},
				},
			},
		).
		WithNearVector(nearVector).
		WithWhere(filter).
		WithLimit(limit).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	get := result.Data["Get"].(map[string]any)
	docs := get["Document"].([]interface{})
	var texts []string
	for _, item := range docs {
		doc := item.(map[string]any)
		text := doc["text"].(string)
		texts = append(texts, text)
	}
	return texts, nil
}

func getAskHandler(client *weaviate.Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(400)
			return
		}

		var r AskRequestBody
		err := json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer req.Body.Close()

		embedding, err := getEmbedding(r.Question)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		results, err := similaritySearch(client, embedding, r.Source, 5)
		if err != nil {
			http.Error(w, fmt.Sprintf("%w", err), 400)
			return
		}

		source := strings.Join(results, "\n\n")
		body := OllamaGenerate{
			Model: "qwen2.5:7b",
			Prompt: fmt.Sprintf(`Answer the question only based on the following context:

      Context: "%w"

      Question: "%w"
      `, source, r.Question),
		}
		data, err := json.Marshal(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("%w", err), 400)
		}

		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}

		resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(data))
		if err != nil {
			http.Error(w, fmt.Sprintf("%w", err), 400)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	}
}
