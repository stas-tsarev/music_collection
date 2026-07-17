package api

import (
	"net/http"
	"pr_1_music_collection/pkg/repository"

	"github.com/gorilla/mux"
)

type api struct {
	r          *mux.Router
	pgdatabase *repository.PGRepository
}

func NewAPI(router *mux.Router, db *repository.PGRepository) *api {
	return &api{r: router, pgdatabase: db}
}

func (api *api) Handles() {
	api.r.HandleFunc("/api/tracks/{id}", api.trackByID).Methods(http.MethodGet)
	api.r.HandleFunc("/api/tracks", api.trackList).Methods(http.MethodGet)
}

func (api *api) ListenAndServe(address string) error {
	return http.ListenAndServe(address, api.r)
}
