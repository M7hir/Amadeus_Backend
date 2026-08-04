package ytDLP

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	ytdlp "github.com/lrstanley/go-ytdlp"
)

var (
	controlCharsRegex = regexp.MustCompile(`[\x00-\x1F\x7F]`)
	allowedCharsRegex = regexp.MustCompile(`[^\p{L}\p{N}\s\-_'",.\(\)&!?]`)
)

func ResolveStream(title, artist string) (string, string, error) {
	url, err := getStreamURL("ytsearch1", title, artist)
	if err == nil && url != "" {
		return url, "youtube", nil
	}

	url, err = getStreamURL("scsearch1", title, artist)
	if err == nil && url != "" {
		return url, "soundcloud", nil
	}

	return "", "", errors.New("no stream found")
}

func getStreamURL(source, title, artist string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := source + ":" + artist + " " + title + " audio"

	dl := ytdlp.New().
		Format("bestaudio").
		NoPlaylist().
		Print("urls")

	result, err := dl.Run(ctx, query)
	if err != nil {
		return "", err
	}

	for _, entry := range result.OutputLogs {
		line := strings.TrimSpace(entry.Line)
		if strings.HasPrefix(line, "http") {
			return line, nil
		}
	}

	return "", errors.New("no URL in output")
}

func SanitizeStreamQuery(input string) (string, error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "", errors.New("empty input")
	}

	if utf8.RuneCountInString(input) > 200 {
		return "", errors.New("input too long")
	}

	input = controlCharsRegex.ReplaceAllString(input, "")
	input = allowedCharsRegex.ReplaceAllString(input, "")

	return strings.TrimSpace(input), nil
}
