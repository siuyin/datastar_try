package main

import (
	"example/public"
	"io"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServerFS(public.Content))
	http.HandleFunc("/hello", helloHandler)

	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}
