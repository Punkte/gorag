package main

import (
	"net/http"

	"gorag/handlers"
	"gorag/services"
)

func main() {
	client, err := services.InitWeaviate()
	if err != nil {
		panic(err)
	}

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/upload", handlers.GetUploadHandler(client))
	http.HandleFunc("/ask", handlers.GetAskHandler(client))
	http.HandleFunc("/documents", handlers.GetDocumentsHandler(client))
	http.ListenAndServe(":4557", nil)
}
