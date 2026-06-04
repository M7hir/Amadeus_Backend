package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"amadeus.m7hir.net/internal/album"
	"amadeus.m7hir.net/internal/artist"
	"amadeus.m7hir.net/internal/reccobeats"
	"amadeus.m7hir.net/internal/songs"
)

func (app *application) trackDetailHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readStringIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	body, err := app.reccobeats.TrackDetail(id)
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackRecommendationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := app.reccobeats.TrackRecommendation(r.URL.Query())
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackAlbumHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readStringIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	body, err := app.reccobeats.TrackAlbum(id)
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackMultipleHandler(w http.ResponseWriter, r *http.Request) {
	body, err := app.reccobeats.TrackMultiple(r.URL.Query())
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackAudioFeaturesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readStringIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	body, err := app.reccobeats.TrackAudioFeatures(id)
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) trackUpstreamErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var validationErr *reccobeats.ValidationError
	if errors.As(err, &validationErr) {
		app.badRequestResponse(w, r, errors.New(validationErr.Message))
		return
	}

	var responseErr *reccobeats.ResponseError
	if errors.As(err, &responseErr) {
		app.errorResponse(w, r, responseErr.StatusCode, responseErr.Message)
		return
	}

	app.serverErrorResponse(w, r, err)
}

func (app *application) writeTrackResponse(w http.ResponseWriter, status int, body []byte) error {
	var payload envelope
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	return app.writeJSON(w, status, payload, nil)
}

func (app *application) DeezerSearchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := app.deezer.DeezerSearch(r.URL.Query())
	if err != nil {
		app.trackUpstreamErrorResponse(w, r, err)
		return
	}

	if err := app.writeTrackResponse(w, http.StatusOK, body); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addSongHandler(w http.ResponseWriter, r *http.Request) {
	liked, err := strconv.ParseBool(r.URL.Query().Get("liked"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var track songs.DeezerTrack
	if err := app.readJSON(w, r, &track); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// app.logger.PrintInfo("track handler:", map[string]interface{}{"track": track})

	if err := app.models.Songs.AddLikedSong(liked, &track); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"liked": liked, "message": "song added"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addArtistHandler(w http.ResponseWriter, r *http.Request) {
	liked, err := strconv.ParseBool(r.URL.Query().Get("liked"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var artist artist.DeezerArtist
	if err := app.readJSON(w, r, &artist); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// app.logger.PrintInfo("track handler:", map[string]interface{}{"track": track})

	if err := app.models.Artist.AddLikedArtist(liked, &artist); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"liked": liked, "message": "artist added"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addAlbumHandler(w http.ResponseWriter, r *http.Request) {
	liked, err := strconv.ParseBool(r.URL.Query().Get("liked"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var album album.DeezerAlbum
	if err := app.readJSON(w, r, &album); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// app.logger.PrintInfo("track handler:", map[string]interface{}{"track": track})

	if err := app.models.Album.AddLikedAlbum(liked, &album); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"liked": liked, "message": "Album added"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
