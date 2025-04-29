package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/index"
	"github.com/SayHelloLexa/searcher/pkg/scan"

	"github.com/gorilla/mux"
)

func idxHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(idx)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(storage)
}

func headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

var (
	storage map[string]crawler.Document
	idx     index.Index
)

func main() {
	r := mux.NewRouter()
	r.Use(headersMiddleware)
	r.HandleFunc("/index", idxHandler).Methods(http.MethodGet)
	r.HandleFunc("/docs", docsHandler).Methods(http.MethodGet)

	idx = index.New()

	// Сканим и в консольку кидаем для вида
	var err error
	storage, err = scan.Scan(&idx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(storage)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
