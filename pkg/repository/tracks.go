package repository

import (
	"context"
	"errors"
	"pr_1_music_collection/pkg/models"
)

func (pgrep *PGRepository) GetTrackByID(trackID int) (models.TrackFull, error) {
	var fulltrack models.TrackFull
	var track models.Track
	tr := pgrep.pgxPool.QueryRow(context.Background(), `
		SELECT * FROM tracks
		WHERE track_id = $1
	`, trackID)

	err := tr.Scan(
		&track.TrackID,
		&track.TrackName,
	)
	if err != nil {
		return models.TrackFull{}, errors.New("not found the track")
	}
	fulltrack.TrackID = track.TrackID
	fulltrack.TrackName = track.TrackName

	var artists []models.Artist
	ar, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    artists.artist_id, artists.artist_name 
		FROM 
		    artists
			JOIN track_artist 
				ON artists.artist_id = track_artist.artist_id
		WHERE track_artist.track_id = $1
	`, trackID)
	for ar.Next() {
		var artist models.Artist
		err = ar.Scan(&artist.ArtistID, &artist.ArtistName)
		if err != nil {
			return models.TrackFull{}, err
		}
		artists = append(artists, artist)
	}
	for i, _ := range artists {
		fulltrack.Artists = append(fulltrack.Artists, artists[i].ArtistName)
	}
	if len(fulltrack.Artists) == 0 {
		fulltrack.Artists = append(fulltrack.Artists, "Неизвестный Исполнитель")
	}

	var albums []models.Album
	al, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    albums.album_id, albums.album_name 
		FROM 
		    albums
			JOIN track_album
				ON albums.album_id = track_album.album_id
		WHERE track_album.track_id = $1
	`, trackID)
	for al.Next() {
		var album models.Album
		err = al.Scan(&album.AlbumID, &album.AlbumName)
		if err != nil {
			return models.TrackFull{}, err
		}
		albums = append(albums, album)
	}
	for i, _ := range albums {
		fulltrack.Albums = append(fulltrack.Albums, albums[i].AlbumName)
	}
	if len(fulltrack.Albums) == 0 {
		fulltrack.Albums = append(fulltrack.Albums, "Неизвестный Альбом")
	}

	return fulltrack, nil
}

func (pgrep *PGRepository) GetTracks() ([]models.TrackFull, error) {
	var result []models.TrackFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT track_id FROM tracks
	`)
	if err != nil {
		return []models.TrackFull{}, errors.New("not found the track")
	}

	for ids.Next() {
		var track_id int
		err = ids.Scan(&track_id)
		if err != nil {
			return []models.TrackFull{}, err
		}
		tmp, err := pgrep.GetTrackByID(track_id)
		if err != nil {
			return []models.TrackFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (pgrep *PGRepository) GetTracksByName(trackName string) ([]models.TrackFull, error) {
	var result []models.TrackFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT track_id FROM tracks
		WHERE LOWER(track_name) = LOWER($1);
	`, trackName)
	if err != nil {
		return []models.TrackFull{}, errors.New("not found the track")
	}

	for ids.Next() {
		var track_id int
		err = ids.Scan(&track_id)
		if err != nil {
			return []models.TrackFull{}, err
		}
		tmp, err := pgrep.GetTrackByID(track_id)
		if err != nil {
			return []models.TrackFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (pgrep *PGRepository) CreateTrack(fulltrack models.TrackFull) error {
	pgrep.mu.Lock()
	defer pgrep.mu.Unlock()

	_, err := pgrep.pgxPool.Exec(context.Background(), `
		INSERT INTO tracks (track_name)
		VALUES ($1);
	`, fulltrack.TrackName)
	return err
}
