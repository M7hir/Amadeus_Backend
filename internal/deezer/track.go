package deezer

import (
	"io"
	"net/http"
	"net/url"
)

type Track struct {
	ID                    string `json:"id"`
	Readable              bool   `json:"readable"`
	Title                 string `json:"title"`
	TitleShort            string `json:"title_short"`
	TitleVersion          string `json:"title_version"`
	ISRC                  string `json:"isrc"`
	Link                  string `json:"link"`
	Duration              int    `json:"duration,string"`
	Rank                  int    `json:"rank,string"`
	ExplicitLyrics        bool   `json:"explicit_lyrics"`
	ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
	ExplicitContentCover  int    `json:"explicit_content_cover"`
	Preview               string `json:"preview"`
	MD5Image              string `json:"md5_image"`
	Artist                Artist `json:"artist"`
	Album                 Album  `json:"album"`
	Type                  string `json:"type"`
}

type Artist struct {
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

type Album struct {
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

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ResponseError struct {
	StatusCode int
	Message    string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (e *ResponseError) Error() string {
	return e.Message
}

func (c *Client) DeezerSearch(q url.Values) ([]byte, error) {
	// logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	// logger.PrintInfo("query:", map[string]interface{}{"query": q})
	return c.DeezerGet(q.Encode())
}

func (c *Client) DeezerGet(rawQuery string) ([]byte, error) {
	deezerURL := "https://api.deezer.com/search/track"

	if rawQuery != "" {
		deezerURL = deezerURL + "?" + rawQuery
	}
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodGet, deezerURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &ResponseError{StatusCode: res.StatusCode, Message: http.StatusText(res.StatusCode)}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
