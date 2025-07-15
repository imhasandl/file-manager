package main

import (
	"log"
	"net/http"
	"os"

	"github.com/imhasandl/file-manager/handlers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Can not load .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Set port in .env file")
	}

	os.Mkdir("archives", 0755)

	apiConfig := handlers.NewAPIConfig()

	mux := http.NewServeMux()

    mux.HandleFunc("POST /tasks", apiConfig.CreateTask)
    mux.HandleFunc("POST /tasks/{taskID}/files", apiConfig.AddFile)
    mux.HandleFunc("GET /tasks/{taskID}/status", apiConfig.GetTaskStatus)
    mux.HandleFunc("GET /download/{filename}", apiConfig.DownloadArchive)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
