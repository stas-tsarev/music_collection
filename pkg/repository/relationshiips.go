package repository

import (
	"context"
	"errors"
	"pr_1_music_collection/pkg/models"
)

func (pgrep *PGRepository) DeleteTrackAlbumRel(trackAlbum models.TrackAlbum) error {
	_, err := pgrep.pgxPool.Exec(context.Background(), `
		DELETE FROM track_album WHERE track_id = $1 and album_id = $2;
	`, trackAlbum.TrackID, trackAlbum.AlbumID)

	if err != nil {
		return errors.New("cannot delete track_album relationship")
	}

	return err
}

func (pgrep *PGRepository) DeleteTrackArtistRel(trackArtist models.TrackArtist) error {
	_, err := pgrep.pgxPool.Exec(context.Background(), `
		DELETE FROM track_artist WHERE track_id = $1 and artist_id = $2;
	`, trackArtist.TrackID, trackArtist.ArtistID)

	if err != nil {
		return errors.New("cannot delete track_artist relationship")
	}

	return err
}

func (pgrep *PGRepository) DeleteAlbumArtistRel(albumArtist models.AlbumArtist) error {
	_, err := pgrep.pgxPool.Exec(context.Background(), `
		DELETE FROM album_artist WHERE album_id = $1 and artist_id = $2;
	`, albumArtist.AlbumID, albumArtist.ArtistID)

	if err != nil {
		return errors.New("cannot delete album_artist relationship")
	}

	return err
}