package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"amadeus.m7hir.net/internal/reccobeats"
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
