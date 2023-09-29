package main

import (
	"log"
	"net/http"
	"search/res"
)

func main() {
	log.Fatal(http.ListenAndServe(":3001", search.Router()))
}
