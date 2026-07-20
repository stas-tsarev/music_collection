package repository

import (
	"context"
	"errors"
	"pr_1_music_collection/pkg/models"
)

func (pgrep *PGRepository) GetArtistByID(ArtistID int) (models.ArtistFull, error) {
	var fullartist models.ArtistFull
	var artist models.Artist
	ar := pgrep.pgxPool.QueryRow(context.Background(), `
		SELECT * FROM artists
		WHERE artist_id = $1
	`, ArtistID)

	err := ar.Scan(
		&artist.ArtistID,
		&artist.ArtistName,
	)
	if err != nil {
		return models.ArtistFull{}, errors.New("not found the track")
	}
	fullartist.ArtistID = artist.ArtistID
	fullartist.ArtistName = artist.ArtistName

	var albums []models.Album
	al, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    albums.album_id, albums.album_name 
		FROM 
		    albums
			JOIN album_artist 
				ON albums.album_id = album_artist.album_id
		WHERE album_artist.artist_id = $1
	`, ArtistID)
	for al.Next() {
		var album models.Album
		err = al.Scan(&album.AlbumID, &album.AlbumName)
		if err != nil {
			return models.ArtistFull{}, err
		}
		albums = append(albums, album)
	}
	for i, _ := range albums {
		fullartist.Albums = append(fullartist.Albums, albums[i].AlbumName)
	}
	if len(fullartist.Albums) == 0 {
		fullartist.Albums = append(fullartist.Albums, "Альбомы отсутствуют")
	}

	var tracks []models.Track
	tr, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT 
		    tracks.track_id, tracks.track_name 
		FROM 
		    tracks
			JOIN track_artist
				ON tracks.track_id = track_artist.track_id
		WHERE track_artist.artist_id = $1
	`, ArtistID)
	for tr.Next() {
		var track models.Track
		err = tr.Scan(&track.TrackID, &track.TrackName)
		if err != nil {
			return models.ArtistFull{}, err
		}
		tracks = append(tracks, track)
	}
	for i, _ := range tracks {
		fullartist.Tracks = append(fullartist.Tracks, tracks[i].TrackName)
	}
	if len(fullartist.Tracks) == 0 {
		fullartist.Tracks = append(fullartist.Tracks, "Композищии отсутствуют")
	}
	return fullartist, nil
}

func (pgrep *PGRepository) GetArtists() ([]models.ArtistFull, error) {
	var result []models.ArtistFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT artist_id FROM artists
	`)
	if err != nil {
		return []models.ArtistFull{}, errors.New("not found the artist")
	}

	for ids.Next() {
		var artist_id int
		err = ids.Scan(&artist_id)
		if err != nil {
			return []models.ArtistFull{}, err
		}
		tmp, err := pgrep.GetArtistByID(artist_id)
		if err != nil {
			return []models.ArtistFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (pgrep *PGRepository) GetArtistsByName(artistName string) ([]models.ArtistFull, error) {
	var result []models.ArtistFull

	ids, err := pgrep.pgxPool.Query(context.Background(), `
		SELECT artist_id FROM artists
		WHERE LOWER(artist_name) = LOWER($1);
	`, artistName)
	if err != nil {
		return []models.ArtistFull{}, errors.New("not found the artist")
	}

	for ids.Next() {
		var artist_id int
		err = ids.Scan(&artist_id)
		if err != nil {
			return []models.ArtistFull{}, err
		}
		tmp, err := pgrep.GetArtistByID(artist_id)
		if err != nil {
			return []models.ArtistFull{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (pgrep *PGRepository) CreateArtist(artist models.Artist) error {
	_, err := pgrep.pgxPool.Exec(context.Background(), `
		INSERT INTO artists(artist_name)
		VALUES ($1)
	`, artist.ArtistName)

	if err != nil {
		return errors.New("cannot insert the artist")
	}

	return nil
}

func (pgrep *PGRepository) AddAlbumToArtist(albumArtist models.AlbumArtist) error {
	req, _ := pgrep.pgxPool.Query(context.Background(), `
		SELECT * FROM album_artist
		WHERE album_id = $1 AND artist_id = $2
	`, albumArtist.AlbumID, albumArtist.ArtistID)

	tmp := 0
	for req.Next() {
		tmp++
	}

	if tmp == 0 {
		_, err := pgrep.pgxPool.Exec(context.Background(), `
		INSERT INTO album_artist(album_id, artist_id)
		VALUES ($1, $2)
	`, albumArtist.AlbumID, albumArtist.ArtistID)

		if err != nil {
			return errors.New("cannot add album to artist")
		}
	}

	req, _ = pgrep.pgxPool.Query(context.Background(), `
		SELECT track_id FROM track_album WHERE album_id = $1;
	`, albumArtist.AlbumID)

	var trackIDS []int
	for req.Next() {
		var trackID int
		err := req.Scan(&trackID)
		if err != nil {
			return errors.New("cannot get track ids from album")
		}
		trackIDS = append(trackIDS, trackID)
	}

	for _, trackID := range trackIDS {
		req, _ = pgrep.pgxPool.Query(context.Background(), `
			SELECT * FROM track_artist
			WHERE track_id = $1 AND artist_id = $2
		`, trackID, albumArtist.ArtistID)

		tmp = 0
		for req.Next() {
			tmp++
		}

		if tmp == 0 {
			_, err := pgrep.pgxPool.Exec(context.Background(), `
				INSERT INTO track_artist(track_id, artist_id)
				VALUES ($1, $2)
			`, trackID, albumArtist.ArtistID)

			if err != nil {
				return errors.New("cannot add track to artist")
			}
		}
	}

	return nil
}

func (pgrep *PGRepository) AddTrackToArtist(trackArtist models.TrackArtist) error {
	req, _ := pgrep.pgxPool.Query(context.Background(), `
			SELECT * FROM track_artist
			WHERE track_id = $1 AND artist_id = $2
		`, trackArtist.TrackID, trackArtist.ArtistID)

	tmp := 0
	for req.Next() {
		tmp++
	}

	if tmp == 0 {
		_, err := pgrep.pgxPool.Exec(context.Background(), `
				INSERT INTO track_artist(track_id, artist_id)
				VALUES ($1, $2)
			`, trackArtist.TrackID, trackArtist.ArtistID)

		if err != nil {
			return errors.New("cannot add track to artist")
		}
	}

	return nil
}
