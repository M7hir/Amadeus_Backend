package main

import (
	"errors"
	"net/http"

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

	if users.ValidateUser(v, newUser); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.InsertUser(newUser)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	// err = user.PasswordHash.Ser
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": newUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
