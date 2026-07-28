package api

import (
	"encoding/json"
	"net/http"
	"pr_1_music_collection/pkg/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *api) artistByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artist_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := api.pgdatabase.GetArtistByID(artist_id)
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

func (api *api) artistList(w http.ResponseWriter, r *http.Request) {
	data, err := api.pgdatabase.GetArtists()
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

func (api *api) artistByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistName := vars["name"]

	data, err := api.pgdatabase.GetArtistsByName(artistName)
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

func (api *api) artistAdd(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist
	err := json.NewDecoder(r.Body).Decode(&artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.CreateArtist(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) artistAddAlbum(w http.ResponseWriter, r *http.Request) {
	var albumArtist models.AlbumArtist
	err := json.NewDecoder(r.Body).Decode(&albumArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.AddAlbumToArtist(albumArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) artistAddTrack(w http.ResponseWriter, r *http.Request) {
	var trackArtist models.TrackArtist
	err := json.NewDecoder(r.Body).Decode(&trackArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.AddTrackToArtist(trackArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) artistDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artist_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.DeleteArtist(artist_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) artistUpdate(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist
	err := json.NewDecoder(r.Body).Decode(&artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.UpdateArtist(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
