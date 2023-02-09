package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to my awesome site!")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting server at port :3030...")
	http.ListenAndServe(":3030", nil)
}
