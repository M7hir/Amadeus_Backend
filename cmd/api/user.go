package main

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"amadeus.m7hir.net/internal/data"
	users "amadeus.m7hir.net/internal/users"
	"amadeus.m7hir.net/internal/validator"
)

func (app *application) userSignUpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newUser := &users.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Activated: false,
	}

	err = newUser.PasswordHash.Set(input.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	users.ValidateUser(v, newUser)
	if len(v.Errors) != 0 {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.InsertUser(newUser)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	token, err := app.models.Token.New(newUser.Id, 3*24*time.Hour, data.ScopeActivation)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	scheme := "http"

	if r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	} else if r.TLS != nil {
		scheme = "https"
	}

	activationURL := url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   "/v1/users/activated",
	}

	activationQuery := activationURL.Query()
	activationQuery.Set("token", token.Plaintext)
	activationURL.RawQuery = activationQuery.Encode()

	data := map[string]interface{}{
		"activationToken": token.Plaintext,
		"activationLink":  activationURL.String(),
		"userID":          newUser.Id,
	}

	app.background(func() {
		err = app.mailer.Send(newUser.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": newUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if r.Method == http.MethodGet {
		input.TokenPlaintext = r.URL.Query().Get("token")
	} else {
		err := app.readJSON(w, r, &input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.User.UpdateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Token.DeleteAllForUser(data.ScopeActivation, user.Id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

}
