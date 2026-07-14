package main

import (
	"context"
	"net/http"

	user "amadeus.m7hir.net/internal/users"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *user.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *user.User {
	user, ok := r.Context().Value(userContextKey).(*user.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
