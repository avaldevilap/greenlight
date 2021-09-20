package main

import (
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		if err := app.mailer.Send(user.Email, "user_welcome.tmpl", data); err != nil {
			app.logger.Printf("Error sending welcome email: %s", err)
		}
	}()

	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
