package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//reccobeats
	router.HandlerFunc(http.MethodGet, "/v1/track-recommendation", app.trackRecommendationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tracks", app.trackMultipleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id/album", app.trackAlbumHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id/audio-features", app.trackAudioFeaturesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/track/:id", app.trackDetailHandler)

	//deezer
	router.HandlerFunc(http.MethodGet, "/v1/search/track", app.DeezerSearchHandler)
	router.HandlerFunc(http.MethodPost, "/v1/song", app.addSongHandler)
	router.HandlerFunc(http.MethodPost, "/v1/artist", app.addArtistHandler)
	router.HandlerFunc(http.MethodPost, "/v1/album", app.addAlbumHandler)

	router.HandlerFunc(http.MethodPost, "/v1/signup", app.userSignUpHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}

