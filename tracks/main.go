package main

import (
	"log"
	"net/http"
	"tracks/repository"
	"tracks/res"
)

func main() {
	repository.Init()
	repository.Create()
	repository.Clear()
	log.Fatal(http.ListenAndServe(":3000", tracks.Router()))
}
