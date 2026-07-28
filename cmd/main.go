package main

import (
	"log"
	"pr_1_music_collection/pkg/api"
	"pr_1_music_collection/pkg/repository"

	"github.com/gorilla/mux"
)

const addr = "localhost:8080"
const conn = "postgres://postgres:s2t0a0s3@localhost:5432/music"

func main() {
	db, err := repository.NewPGRepository(conn)
	if err != nil {
		log.Fatal(err)
	}

	api := api.NewAPI(mux.NewRouter(), db)
	api.Handles()
	log.Fatal(api.ListenAndServe(addr))
}
