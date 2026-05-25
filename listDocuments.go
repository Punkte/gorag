package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

func getDocumentsHandler(client *weaviate.Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		groupedBy := graphql.Field{
			Name: "groupedBy", Fields: []graphql.Field{
				{Name: "value"},
			},
		}

		result, err := client.GraphQL().Aggregate().
			WithClassName("Document").
			WithGroupBy("source").
			WithFields(groupedBy).
			Do(context.Background())

		if err != nil {
			w.WriteHeader(400)
			return
		}

		aggregate := result.Data["Aggregate"].(map[string]any)
		docs := aggregate["Document"].([]interface{})

		var sources []string
		for _, item := range docs {
			doc := item.(map[string]any)
			groupedByField := doc["groupedBy"].(map[string]any)
			value := groupedByField["value"].(string)
			sources = append(sources, value)
		}
		data, err := json.Marshal(sources)
		if err != nil {
			panic(fmt.Errorf("an error occurred: %w", err))
		}
		w.Write(data)
	}
}
