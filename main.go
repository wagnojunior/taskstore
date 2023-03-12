// This project is largely based on Eli Bendersky's series `REST Servers in Go`.
// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

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
	mux.HandleFunc("/tag/", server.TagHandler)
	mux.HandleFunc("/due/", server.DueHandler)

	fmt.Println("Starting server at port :3030...")
	log.Fatal(http.ListenAndServe("localhost:3030", mux))
}
