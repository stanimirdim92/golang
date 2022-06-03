package main

import (
	"flag"
	"log"
	"net/http"
	"urlshortener/handlers"
)

var (
	//shortPathsJSON = flag.String("shortPathsJSON", "shortPaths.json", "The file containing shortened paths to URL's")
	//shortPathsYAML = flag.String("shortPathsYAML", "shortPaths.yaml", "The file containing shortened paths to URL's")
	shortPathsMAP = map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
)

type shortenerHandler struct {
	url string
}

func (sh shortenerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(sh.url))
	if err != nil {
		return
	}
}

func main() {
  	mux := defaultMux()
	mux.Handle("/shorten", shortenerHandler{url:"google.com"})
	urlshortener.MapHandler(shortPathsMAP, mux)
	log.Fatal(http.ListenAndServe(":8082", mux))
}


func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}
