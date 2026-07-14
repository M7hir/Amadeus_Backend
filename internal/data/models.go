package data

import (
	"database/sql"

	"amadeus.m7hir.net/internal/album"
	"amadeus.m7hir.net/internal/artist"
	"amadeus.m7hir.net/internal/jsonlog"
	"amadeus.m7hir.net/internal/songs"
	user "amadeus.m7hir.net/internal/users"
)

var (
	ErrRecordNotFound = user.ErrRecordNotFound
	ErrEditConflict   = user.ErrEditConflict
)

type Models struct {
	Songs  songs.SongsModel
	Artist artist.ArtistModel
	Album  album.AlbumModel
	User   user.UserModel
	Token  TokenModel
}

func NewModels(db *sql.DB, logger *jsonlog.Logger) Models {
	return Models{
		Songs:  songs.SongsModel{DB: db, Logger: logger},
		Artist: artist.ArtistModel{DB: db, Logger: logger},
		Album:  album.AlbumModel{DB: db, Logger: logger},
		User:   user.UserModel{DB: db, Logger: logger},
		Token:  TokenModel{DB: db},
	}
}
