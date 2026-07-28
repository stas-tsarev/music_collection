package main

import (
	"log"
	"os"
	"pr_1_music_collection/pkg/api"
	"pr_1_music_collection/pkg/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")
	conn := os.Getenv("POSTGR")

	db, err := repository.NewPGRepository(conn)
	if err != nil {
		log.Fatal(err)
	}

	api := api.NewAPI(mux.NewRouter(), db)
	api.Handles()
	log.Fatal(api.ListenAndServe(addr))
}
