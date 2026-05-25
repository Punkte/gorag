package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorag/services"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func saveFile(req *http.Request) (string, error) {
	req.ParseMultipartForm(32 << 20)
	file, fileHeader, err := req.FormFile("file")

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer file.Close()

	filename := fileHeader.Filename
	h := sha256.New()
	h.Write([]byte(time.Now().String()))
	filename = fmt.Sprintf("./files/%s_%s", hex.EncodeToString(h.Sum(nil)), filename)
	createdFile, createErr := os.Create(filename)

	if createErr != nil {
		return "", fmt.Errorf("%w", createErr)
	}
	defer createdFile.Close()

	_, copyErr := io.Copy(createdFile, file)
	if copyErr != nil {
		return "", fmt.Errorf("%w", copyErr)
	}

	return filename, nil
}

func GetUploadHandler(client *weaviate.Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(400)
			return
		}

		filename, err := saveFile(req)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer req.Body.Close()

		go func() {
			content, err := services.ExtractText(filename)
			if err != nil {
				log.Printf("extractText error: %v", err)
			}
			chunks := services.ChunkText(content, 800, 150)

			for i := range chunks {
				embedding, err := services.GetEmbedding(chunks[i])
				if err != nil {
					log.Printf("embedding error: %v", err)
				}
				err = services.StoreChunk(client, chunks[i], filename, embedding)
				if err != nil {
					log.Printf("error while storing: %v", err)
				}
			}
		}()
	}
}
