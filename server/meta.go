package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func metaURLsHandler(urls []string) urlHandler {
	res := urlHandler{
		name:     "API URLs",
		url:      "/urls",
		help:     "Lists all API urls.",
		examples: []string{"/urls"},
		handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, "Served URLS:\n\n")
			fmt.Fprint(w, strings.Join(urls, "\n"))
		},
	}
	return res
}

// JSONURLExample : JSON container
type JSONURLExample struct {
	Template string `json:"template"`
	URL      string `json:"url"`
}

var metaExamplesHandler = urlHandler{
	name:     "API URL examples",
	url:      "/examples",
	help:     "Lists all API urls examples.",
	examples: []string{"/examples"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		res := []JSONURLExample{}
		for _, subRouter := range subRouters {
			for _, handler := range subRouter.handlers {
				for _, example := range handler.examples {
					//url := "http://localhost:" + port + subRouter.root + example
					url := subRouter.root + example
					template := subRouter.root + handler.url
					ex := JSONURLExample{Template: template, URL: url}
					res = append(res, ex)
				}
			}
		}
		js, err := json.Marshal(res)
		if err != nil {
			log.Printf("lexserver: failed to marshal struct : %v", err)
			http.Error(w, fmt.Sprintf("failed to marshal struct : %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(js))
	},
}
