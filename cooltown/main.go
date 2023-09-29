package main

import (
	"cooltown/res"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":3002", cooltown.Router()))
}
