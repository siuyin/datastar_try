package main

import (
	"encoding/json"
	"example/public"
	"fmt"
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
	http.HandleFunc("/boilwater", boilWaterHandler)

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

func boilWaterHandler(w http.ResponseWriter, r *http.Request) {
	type Signals struct {
		Tmp json.Number `json:"tmp"`
	}
	sig := &Signals{}
	if err := datastar.ReadSignals(r, sig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Println(sig.Tmp)

	sse := datastar.NewSSE(w, r)
	startTmp, err := sig.Tmp.Float64()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	sse.MergeFragments(`<div id="tmpincr">Changing temperature: ` + fmt.Sprintf("%f with %f = 212°F", startTmp, 212.0-startTmp))
	sse.MergeSignals([]byte("{tmp: 212}"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World\n")
}
