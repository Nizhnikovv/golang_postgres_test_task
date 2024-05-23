package main

import "log"

func main() {
	storage, err := NewPostgresStorage()
	if err != nil {
		log.Fatalf("Error while creating storage: %v", err)
	}

	if err := storage.Init(); err != nil {
		log.Fatalf("Error while initializing storage: %v", err)
	}

	server := NewAPIServer(":3246", storage)
	if err := server.Run(); err != nil {
		log.Fatalf("Error while running server: %v", err)
	}
}
