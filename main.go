package main

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar/sdk/go"
)

//go:embed hello-world.html
var helloWorldHTML []byte

func main() {
	r := chi.NewRouter()

	const message = "你好, Datastar!"
	type Store struct {
		Delay time.Duration `json:"delay"` // delay in milliseconds between each character of the message.
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(helloWorldHTML)
	})

	r.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)
		runes := []rune(message)  // 转换为rune数组
		
		for i := 0; i < len(runes); i++ {
			sse.MergeFragments(`<div id="message">` + string(runes[:i+1]) + `</div>`)
			time.Sleep(store.Delay * time.Millisecond)
		}
	})

	http.ListenAndServe(":8080", r)
}
