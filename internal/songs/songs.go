package songs

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type SongsModel struct {
	DB *sql.DB
}

type Song struct {
	id                    string
	title                 string
	artist_name           string
	duration_ms           int
	cover_image_url       string
	created_at            time.Time
	updated_at            time.Time
	artist_id             string
	album_id              string
	explicit_lyrics       bool
	isrc                  string
	rank                  int
	link                  string
	expicit_content_cover int
}

func ValidateSong(song *Song) error {
	if song == nil {
		return errors.New("song cannot be nil")
	}

	if strings.TrimSpace(song.id) == "" {
		return errors.New("song id is required")
	}

	if strings.TrimSpace(song.title) == "" {
		return errors.New("song title is required")
	}

	if strings.TrimSpace(song.artist_name) == "" {
		return errors.New("song artist_name is required")
	}

	if strings.TrimSpace(song.artist_id) == "" {
		return errors.New("song artist_id is required")
	}

	if strings.TrimSpace(song.album_id) == "" {
		return errors.New("song album_id is required")
	}

	if song.duration_ms <= 0 {
		return errors.New("song duration_ms must be greater than 0")
	}

	if song.rank < 0 {
		return errors.New("song rank cannot be negative")
	}

	if song.expicit_content_cover < 0 {
		return errors.New("song expicit_content_cover cannot be negative")
	}

	return nil
}

type DeezerTrack struct {
	ID                    string       `json:"id"`
	Readable              bool         `json:"readable"`
	Title                 string       `json:"title"`
	TitleShort            string       `json:"title_short"`
	TitleVersion          string       `json:"title_version"`
	ISRC                  string       `json:"isrc"`
	Link                  string       `json:"link"`
	Duration              int          `json:"duration,string"`
	Rank                  int          `json:"rank,string"`
	ExplicitLyrics        bool         `json:"explicit_lyrics"`
	ExplicitContentLyrics int          `json:"explicit_content_lyrics"`
	ExplicitContentCover  int          `json:"explicit_content_cover"`
	Preview               string       `json:"preview"`
	MD5Image              string       `json:"md5_image"`
	Artist                DeezerArtist `json:"artist"`
	Album                 DeezerAlbum  `json:"album"`
	Type                  string       `json:"type"`
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
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}

type DeezerAlbum struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXL     string `json:"cover_xl"`
	MD5Image    string `json:"md5_image"`
	Tracklist   string `json:"tracklist"`
	Type        string `json:"type"`
}

func TrackToSong(track *DeezerTrack) *Song {
	now := time.Now()
	return &Song{
		id:                    track.ID,
		title:                 track.Title,
		artist_name:           track.Artist.Name,
		duration_ms:           track.Duration,
		cover_image_url:       track.Album.Cover,
		created_at:            now,
		updated_at:            now,
		artist_id:             track.Artist.ID,
		album_id:              track.Album.ID,
		explicit_lyrics:       track.ExplicitLyrics,
		isrc:                  track.ISRC,
		rank:                  track.Rank,
		link:                  track.Link,
		expicit_content_cover: track.ExplicitContentCover,
	}
}

func (s SongsModel) InsertLikedSong(liked bool, track *Song) error {
	if err := ValidateSong(track); err != nil {
		return err
	}

	query := `INSERT INTO songs (
		id, title, artist_name, duration_ms, cover_image_url, 
		created_at, updated_at, artist_id, album_id, 
		explicit_lyrics, isrc, rank, link,expicit_content_cover
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14
	) ON CONFLICT (id) DO NOTHING`

	result, err := s.DB.Exec(
		query,
		track.id,
		track.title,
		track.artist_name,
		track.duration_ms,
		track.cover_image_url,
		track.created_at,
		track.updated_at,
		track.artist_id,
		track.album_id,
		track.explicit_lyrics,
		track.isrc,
		track.rank,
		track.link,
		track.expicit_content_cover,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("song already exists")
	}

	return nil
}

func (s SongsModel) AddLikedSong(liked bool, track *DeezerTrack) error {
	song := TrackToSong(track)
	return s.InsertLikedSong(liked, song)
}
