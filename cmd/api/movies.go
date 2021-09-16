package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/avaldevilap/greenlight/internal/data"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	validate = validator.New()

	if err = validate.Struct(movie); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Read the JSON data from the request body and store it in the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update the movie record with the data from the input struct
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	// Validate the movie record
	validate = validator.New()
	if err = validate.Struct(movie); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	// Update the movie record in the database
	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return a 201 Created response with the movie record in the response body
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
