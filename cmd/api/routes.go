package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track-recommendation", app.trackRecommendationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tracks", app.trackMultipleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id/album", app.trackAlbumHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id/audio-features", app.trackAudioFeaturesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id", app.trackDetailHandler)

	return router
}
