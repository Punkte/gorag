package main

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func initWeaviate() (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	existing, _ := client.Schema().ClassGetter().WithClassName("Document").Do(context.Background())
	if existing != nil {
		return client, nil
	}

	documentClass := &models.Class{
		Class:       "Document",
		Description: "A collection of documents",
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: schema.DataTypeText.PropString(),
			},
			{
				Name:     "source",
				DataType: schema.DataTypeText.PropString(),
			},
		},
	}

	err = client.Schema().ClassCreator().WithClass(documentClass).Do(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("Document collection created successfully")
	return client, nil
}

func storeChunk(client *weaviate.Client, text string, source string, embedding []float32) error {
	_, err := client.Data().Creator().
		WithClassName("Document").
		WithProperties(map[string]interface{}{
			"text":   text,
			"source": source,
		}).
		WithVector(embedding).
		Do(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("Chunk stored successfully")
	return nil
}
