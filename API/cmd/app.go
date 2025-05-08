package main

import (
	"github.com/SayHelloLexa/searcher/API"
	"net/http"
)

func main() {
	api := API.New()
	http.ListenAndServe(":8080", api.Router)
}
