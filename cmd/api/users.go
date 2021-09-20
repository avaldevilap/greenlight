package main

import (
	"net/http"

	"github.com/avaldevilap/greenlight/internal/data"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		Activated: false,
	}

	if err := user.Validate(); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	if err := app.models.Users.Insert(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.mailer.Send(user.Email, "user_welcome.tmpl", user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
