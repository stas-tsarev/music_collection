package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *api) albumByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	album_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := api.pgdatabase.GetAlbumByID(album_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) albumList(w http.ResponseWriter, r *http.Request) {
	data, err := api.pgdatabase.GetAlbums()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) albumByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumName := vars["name"]

	data, err := api.pgdatabase.GetAlbumsByName(albumName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
