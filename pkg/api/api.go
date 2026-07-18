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
	api.r.HandleFunc("/api/tracks", api.trackByID).Queries("id", "{id}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/tracks", api.trackByName).Queries("name", "{name}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/tracks", api.trackList).Methods(http.MethodGet)

	api.r.HandleFunc("/api/albums", api.albumByID).Queries("id", "{id}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/albums", api.albumByName).Queries("name", "{name}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/albums", api.albumList).Methods(http.MethodGet)

	api.r.HandleFunc("/api/artists", api.artistByID).Queries("id", "{id}").Methods(http.MethodGet)
}

func (api *api) ListenAndServe(address string) error {
	return http.ListenAndServe(address, api.r)
}
