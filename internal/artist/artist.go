package artist

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"amadeus.m7hir.net/internal/jsonlog"
)

type ArtistModel struct {
	DB     *sql.DB
	Logger *jsonlog.Logger
}

type Artist struct {
	id             string
	picture        string
	picture_big    string
	picture_small  string
	picture_medium string
	picture_xl     string
	tracklist      string
	link           string
	title          string
	created_at     time.Time
	updated_at     time.Time
}

type DeezerArtist struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXL     string `json:"picture_xl"`
	NbAlbum       int    `json:"nb_album"`
	NbFan         int    `json:"nb_fan"`
	Radio         bool   `json:"radio"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}

func DeezerArtistToLocal(artist *DeezerArtist) *Artist {
	now := time.Now()
	return &Artist{
		id:             artist.ID,
		picture:        artist.Picture,
		picture_big:    artist.PictureBig,
		picture_small:  artist.PictureSmall,
		picture_medium: artist.PictureMedium,
		picture_xl:     artist.PictureXL,
		tracklist:      artist.Tracklist,
		link:           artist.Link,
		title:          artist.Name,
		created_at:     now,
		updated_at:     now,
	}
}

func ValidateArtist(artist *Artist) error {
	if artist == nil {
		return errors.New("artist cannot be nil")
	}
	if strings.TrimSpace(artist.id) == "" {
		return errors.New("artist id is required")
	}

	if strings.TrimSpace(artist.title) == "" {
		return errors.New("artist title is required")
	}
	return nil
}

func (a ArtistModel) InsertLikedArtist(liked bool, artist *Artist) error {
	if err := ValidateArtist(artist); err != nil {
		return err
	}

	query := `INSERT INTO artist (
		id,picture,picture_big,picture_small,picture_medium,
		picture_xl,tracklist,link,title,created_at,updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	 ON CONFLICT (id) DO NOTHING`

	result, err := a.DB.Exec(
		query,
		artist.id,
		artist.picture,
		artist.picture_big,
		artist.picture_small,
		artist.picture_medium,
		artist.picture_xl,
		artist.tracklist,
		artist.link,
		artist.title,
		artist.created_at,
		artist.updated_at,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("artist already exists")
	}
	return nil
}

func (a ArtistModel) AddLikedArtist(liked bool, artist *DeezerArtist) error {
	newArtist := DeezerArtistToLocal(artist)
	return a.InsertLikedArtist(liked, newArtist)
}
