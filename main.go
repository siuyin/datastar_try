package main

import (
	"example/public"
	"io"
	"log"
	"net/http"
	"time"

	datastar "github.com/starfederation/datastar/sdk/go"
)

func main() {
	http.Handle("/", http.FileServerFS(public.Content))
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/hello", helloHandler)

	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func timeHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	for n := 0; n < 500; n++ {
		sse.MergeFragments(`<div id="mytime">` + time.Now().Format("15:04:05.0000") + `</div>`)
		time.Sleep(30 * time.Millisecond)
	}
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}
