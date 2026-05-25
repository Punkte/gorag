package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type OllamaEmbedBody struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type OllamaEmbedResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   int         `json:"total_duration"`
	LoadDuration    int         `json:"load_duration"`
	PromptEvalCount int         `json:"prompt_eval_count"`
}

func getEmbedding(text string) ([]float32, error) {
	body := OllamaEmbedBody{
		Model: "nomic-embed-text",
		Input: text,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return []float32{}, err
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	resp, err := http.Post(ollamaURL+"/api/embed", "application/json", bytes.NewReader(data))
	if err != nil {
		return []float32{}, err
	}
	defer resp.Body.Close()

	var r OllamaEmbedResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return []float32{}, err
	}

	return r.Embeddings[0], nil
}
