package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *api) trackByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	track_id, err := strconv.Atoi(vars["id"])

	//fmt.Println(1)

	data, err := api.pgdatabase.GetTrackByID(track_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(2)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println(3)
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
