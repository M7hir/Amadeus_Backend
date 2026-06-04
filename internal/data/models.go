package data

import (
	"database/sql"
	"errors"

	"amadeus.m7hir.net/internal/album"
	"amadeus.m7hir.net/internal/artist"
	"amadeus.m7hir.net/internal/songs"
)

var (
	ErrRecordNotFound = errors.New("Record Not Found")
)

type Models struct {
	Songs  songs.SongsModel
	Artist artist.ArtistModel
	Album  album.AlbumModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Songs:  songs.SongsModel{DB: db},
		Artist: artist.ArtistModel{DB: db},
		Album:  album.AlbumModel{DB: db},
	}
}
