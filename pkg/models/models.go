package models

type Track struct {
	TrackID   int    `json:"track_id"`
	TrackName string `json:"track_name"`
}

type Album struct {
	AlbumID   int    `json:"album_id"`
	AlbumName string `json:"album_name"`
	Release   int    `json:"release_year"`
}

type Artist struct {
	ArtistID   int    `json:"artist_id"`
	ArtistName string `json:"artist_name"`
}

type TrackFull struct {
	TrackID   int      `json:"track_id"`
	TrackName string   `json:"track_name"`
	Artists   []string `json:"artists_name"`
	Albums    []string `json:"albums_name"`
}
