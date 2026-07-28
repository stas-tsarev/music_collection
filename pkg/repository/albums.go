package repository

import (
	"context"
	"errors"
	"pr_1_music_collection/pkg/models"
	"time"
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

func (pgrep *PGRepository) CreateAlbum(album models.Album) error {
	releaseYear := time.Now().Year()
	if album.Release != 0 {
		releaseYear = album.Release
	}

	_, err := pgrep.pgxPool.Exec(context.Background(), `
		INSERT INTO albums(album_name, release_year)
		VALUES ($1, $2)
	`, album.AlbumName, releaseYear)

	if err != nil {
		return errors.New("cannot insert the album")
	}

	return nil
}

func (pgrep *PGRepository) AddTrackInAlbum(trackAlbum models.TrackAlbum) error {
	req, _ := pgrep.pgxPool.Query(context.Background(), `
		SELECT * FROM track_album
		WHERE track_id = $1 AND album_id = $2
	`, trackAlbum.TrackID, trackAlbum.AlbumID)

	tmp := 0
	for req.Next() {
		tmp++
	}

	if tmp == 0 {
		_, err := pgrep.pgxPool.Exec(context.Background(), `
		INSERT INTO track_album(track_id, album_id)
		VALUES ($1, $2)
	`, trackAlbum.TrackID, trackAlbum.AlbumID)

		if err != nil {
			return errors.New("cannot add track in album")
		}
	}
	return nil
}

func (pgrep *PGRepository) DeleteAlbum(albumID int) error {
	tx, err := pgrep.pgxPool.Begin(context.Background())
	if err != nil {
		return errors.New("cannot begin transaction")
	}
	defer tx.Rollback(context.Background())

	commands := []string{
		"DELETE FROM albums WHERE album_id = $1",
		"DELETE FROM track_album WHERE album_id = $1",
		"DELETE FROM album_artist WHERE album_id = $1",
	}

	for _, command := range commands {
		if _, err = tx.Exec(context.Background(), command, albumID); err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (pgrep *PGRepository) UpdateAlbum(album models.Album) error {
	pgrep.mu.Lock()
	defer pgrep.mu.Unlock()

	var oldAlbum models.Album
	al := pgrep.pgxPool.QueryRow(context.Background(), `
		SELECT * FROM albums WHERE album_id = $1;
	`, album.AlbumID)

	err := al.Scan(&oldAlbum.AlbumID, &oldAlbum.AlbumName, &oldAlbum.Release)
	if err != nil {
		return err
	}

	if oldAlbum.Release <= 0 {
		oldAlbum.Release = time.Now().Year()
	}

	if oldAlbum.AlbumName == album.AlbumName && oldAlbum.Release == album.Release {
		return nil
	}

	_, err = pgrep.pgxPool.Exec(context.Background(), `
		UPDATE albums SET album_name = $1, release_year = $2
		WHERE album_id = $3;
	`, album.AlbumName, album.Release, album.AlbumID)

	return err
}
