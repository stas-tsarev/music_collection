package api

import (
	"encoding/json"
	"net/http"
	"pr_1_music_collection/pkg/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *api) trackByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	track_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := api.pgdatabase.GetTrackByID(track_id)
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

func (api *api) trackList(w http.ResponseWriter, r *http.Request) {
	data, err := api.pgdatabase.GetTracks()
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

func (api *api) trackByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackName := vars["name"]

	data, err := api.pgdatabase.GetTracksByName(trackName)
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

func (api *api) trackAdd(w http.ResponseWriter, r *http.Request) {
	var track models.Track
	err := json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.CreateTrack(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) trackDel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	track_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.DeleteTrack(track_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) trackUpdate(w http.ResponseWriter, r *http.Request) {
	var track models.Track
	err := json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.UpdateTrack(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
