package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, req *http.Request) {
	m := make(map[string]string)
	m["status"] = "ok"
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("an error occurred: %w", err))
	}
	defer req.Body.Close()
	w.Write(data)
}
