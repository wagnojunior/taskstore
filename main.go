package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wagnojunior/taskstore/internal/taskstore"
)

func main() {
	// Starts the taskstore server
	server := taskstore.NewTaskServer()

	// Define mux and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/task/", server.TaskHandler)

	fmt.Println("Starting server at port :3030...")
	log.Fatal(http.ListenAndServe("localhost:3030", mux))
}
