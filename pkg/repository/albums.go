package repository

import (
	"context"
	"errors"
	"pr_1_music_collection/pkg/models"
)

func (pgrep *PGRepository) GetAlbumByID(AlbumID int) (models.AlbumFull, error) {
	var fullalbum models.AlbumFull
	var album models.Album
	al := pgrep.pgxPool.QueryRow(context.Background(), `
		SELECT * FROM albums
		WHERE album_id = $1
	`, AlbumID)

	err := al.Scan(
		&album.AlbumID,
		&album.AlbumName,
		&album.Release,
	)
	if err != nil {
		return models.AlbumFull{}, errors.New("not found the album")
	}
	fullalbum.AlbumID = album.AlbumID
	fullalbum.AlbumName = album.AlbumName
	fullalbum.Release = album.Release

	var artists []models.Artist
	ar, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    artists.artist_id, artists.artist_name 
		FROM 
		    artists
			JOIN album_artist 
				ON artists.artist_id = album_artist.artist_id
		WHERE album_artist.album_id = $1
	`, AlbumID)
	for ar.Next() {
		var artist models.Artist
		err = ar.Scan(&artist.ArtistID, &artist.ArtistName)
		if err != nil {
			return models.AlbumFull{}, err
		}
		artists = append(artists, artist)
	}
	for i, _ := range artists {
		fullalbum.Artists = append(fullalbum.Artists, artists[i].ArtistName)
	}
	if len(fullalbum.Artists) == 0 {
		fullalbum.Artists = append(fullalbum.Artists, "Неизвестный Исполнитель")
	}

	var tracks []models.Track
	tr, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    tracks.track_id, tracks.track_name 
		FROM 
		    tracks
			JOIN track_album
				ON tracks.track_id = track_album.track_id
		WHERE track_album.album_id = $1
	`, AlbumID)
	for tr.Next() {
		var track models.Track
		err = tr.Scan(&track.TrackID, &track.TrackName)
		if err != nil {
			return models.AlbumFull{}, err
		}
		tracks = append(tracks, track)
	}
	for i, _ := range tracks {
		fullalbum.Tracks = append(fullalbum.Tracks, tracks[i].TrackName)
	}
	if len(fullalbum.Tracks) == 0 {
		fullalbum.Tracks = append(fullalbum.Tracks, "Композищии отсутствуют")
	}

	return fullalbum, nil
}

func (pgrep *PGRepository) GetAlbums() ([]models.AlbumFull, error) {
	var result []models.AlbumFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT album_id FROM albums
	`)
	if err != nil {
		return []models.AlbumFull{}, errors.New("not found the album")
	}

	for ids.Next() {
		var album_id int
		err = ids.Scan(&album_id)
		if err != nil {
			return []models.AlbumFull{}, err
		}
		tmp, err := pgrep.GetAlbumByID(album_id)
		if err != nil {
			return []models.AlbumFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (pgrep *PGRepository) GetAlbumsByName(albumName string) ([]models.AlbumFull, error) {
	var result []models.AlbumFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT album_id FROM albums
		WHERE LOWER(album_name) = LOWER($1);
	`, albumName)
	if err != nil {
		return []models.AlbumFull{}, errors.New("not found the album")
	}

	for ids.Next() {
		var album_id int
		err = ids.Scan(&album_id)
		if err != nil {
			return []models.AlbumFull{}, err
		}
		tmp, err := pgrep.GetAlbumByID(album_id)
		if err != nil {
			return []models.AlbumFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}
