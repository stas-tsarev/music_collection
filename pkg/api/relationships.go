package api

import (
	"encoding/json"
	"net/http"
	"pr_1_music_collection/pkg/models"
)

func (api *api) trackAlbumDel(w http.ResponseWriter, r *http.Request) {
	var trackAlbum models.TrackAlbum
	err := json.NewDecoder(r.Body).Decode(&trackAlbum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.DeleteTrackAlbumRel(trackAlbum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) albumArtistDel(w http.ResponseWriter, r *http.Request) {
	var albumArtist models.AlbumArtist
	err := json.NewDecoder(r.Body).Decode(&albumArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.DeleteAlbumArtistRel(albumArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *api) trackArtistDel(w http.ResponseWriter, r *http.Request) {
	var trackArtist models.TrackArtist
	err := json.NewDecoder(r.Body).Decode(&trackArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.pgdatabase.DeleteTrackArtistRel(trackArtist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
