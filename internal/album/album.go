package album

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type AlbumModel struct {
	DB *sql.DB
}

type DeezerAlbum struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Link        string       `json:"link"`
	Cover       string       `json:"cover"`
	CoverSmall  string       `json:"cover_small"`
	CoverMedium string       `json:"cover_medium"`
	CoverBig    string       `json:"cover_big"`
	CoverXL     string       `json:"cover_xl"`
	MD5Image    string       `json:"md5_image"`
	GenreID     int          `json:"genre_id"`
	NbTracks    int          `json:"nb_tracks"`
	RecordType  string       `json:"record_type"`
	Tracklist   string       `json:"tracklist"`
	Explicit    bool         `json:"explicit_lyrics"`
	Artist      DeezerArtist `json:"artist"`
	Type        string       `json:"type"`
}

type DeezerArtist struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Picture      string `json:"picture"`
	PictureSmall string `json:"picture_small"`
	PictureMed   string `json:"picture_medium"`
	PictureBig   string `json:"picture_big"`
	PictureXL    string `json:"picture_xl"`
	Tracklist    string `json:"tracklist"`
	Type         string `json:"type"`
}

type Album struct {
	id           string
	name         string
	cover        string
	cover_big    string
	cover_medium string
	cover_small  string
	cover_xl     string
	link         string
	tracklist    string
	genre_id     int
	created_at   time.Time
	updated_at   time.Time
}

func DeezerAlbumToLocal(album *DeezerAlbum) *Album {
	now := time.Now()
	return &Album{
		id:           album.ID,
		name:         album.Title,
		cover:        album.Cover,
		cover_big:    album.CoverBig,
		cover_medium: album.CoverMedium,
		cover_small:  album.CoverSmall,
		cover_xl:     album.CoverXL,
		link:         album.Link,
		tracklist:    album.Tracklist,
		genre_id:     album.GenreID,
		created_at:   now,
		updated_at:   now,
	}
}

func ValidateAlbum(album *Album) error {
	if album == nil {
		return errors.New("album cannot be nil")
	}
	if strings.TrimSpace(album.id) == "" {
		return errors.New("album id is required")
	}

	if strings.TrimSpace(album.name) == "" {
		return errors.New("album title is required")
	}
	if album.genre_id <= 0 {
		return errors.New("album genre id cannot be zero or less")
	}
	return nil
}

func (a AlbumModel) InsertLikedAlbum(liked bool, album *Album) error {
	if err := ValidateAlbum(album); err != nil {
		return err
	}

	query := `INSERT INTO album (
	id,name,cover,cover_big,cover_medium,cover_small,cover_xl,
link,tracklist,genre_id,created_at,updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	 ON CONFLICT (id) DO NOTHING`

	result, err := a.DB.Exec(
		query,
		album.id,
		album.name,
		album.cover,
		album.cover_big,
		album.cover_medium,
		album.cover_small,
		album.cover_xl,
		album.link,
		album.tracklist,
		album.genre_id,
		album.created_at,
		album.updated_at,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("album already exists")
	}
	return nil
}

func (a AlbumModel) AddLikedAlbum(liked bool, album *DeezerAlbum) error {
	newAlbum := DeezerAlbumToLocal(album)
	return a.InsertLikedAlbum(liked, newAlbum)
}
