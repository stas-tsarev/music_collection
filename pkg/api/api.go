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
	api.r.HandleFunc("/api/tracks", api.trackAdd).Methods(http.MethodPost)

	api.r.HandleFunc("/api/albums", api.albumByID).Queries("id", "{id}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/albums", api.albumByName).Queries("name", "{name}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/albums", api.albumList).Methods(http.MethodGet)
	api.r.HandleFunc("/api/albums", api.albumAdd).Methods(http.MethodPost)
	api.r.HandleFunc("/api/albums/add_track", api.albumAddTrack).Methods(http.MethodPost)

	api.r.HandleFunc("/api/artists", api.artistByID).Queries("id", "{id}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/artists", api.artistByName).Queries("name", "{name}").Methods(http.MethodGet)
	api.r.HandleFunc("/api/artists", api.artistList).Methods(http.MethodGet)
	api.r.HandleFunc("/api/artists", api.artistAdd).Methods(http.MethodPost)
	api.r.HandleFunc("/api/artists/add_track", api.artistAddTrack).Methods(http.MethodPost)
	api.r.HandleFunc("/api/artists/add_album", api.artistAddAlbum).Methods(http.MethodPost)
}

func (api *api) ListenAndServe(address string) error {
	return http.ListenAndServe(address, api.r)
}
