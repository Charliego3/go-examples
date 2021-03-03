package main

import (
	"flag"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"strings"
)

var gz = flag.Bool("gzip", false, "enable automatic gzip compression")

func main() {
	flag.Parse()
	h := wasmContentTypeSetter(http.FileServer(http.Dir("/Users/nzlong/dev/project/Go/temp/WebAssembly/static")))
	if *gz {
		h = gziphandler.GzipHandler(h)
	}

	log.Print("Serving on http://localhost:9999")
	err := http.ListenAndServe(":9999", h)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// http://localhost:9999/index.html
func wasmContentTypeSetter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		h.ServeHTTP(w, r)
	})
}
