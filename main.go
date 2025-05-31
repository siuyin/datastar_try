package main

import (
	"encoding/json"
	"example/public"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/siuyin/dflt"
	datastar "github.com/starfederation/datastar/sdk/go"
)

var (
	baseURL = "/"
	tmpl    *template.Template
)

func init() {
	baseURL = dflt.EnvString("BASE_URL", "/")
	tmpl = template.Must(template.ParseFS(public.Content, "*.html"))
}

type base struct {
	URL string
}

func main() {
	http.Handle("/", http.FileServerFS(public.Content))
	http.HandleFunc("/{$}", indexHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/boilwater", boilWaterHandler)

	log.Println("Starting HTTP server...")
	port := dflt.EnvString("PORT", "8080")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "index.html", base{URL: baseURL}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	for n := 0; n < 50; n++ {
		sse.MergeFragments(`<div id="mytime">` + time.Now().Format("15:04:05.0000") + `</div>`)
		time.Sleep(300 * time.Millisecond)
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
