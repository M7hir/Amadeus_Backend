package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (app *application) albumHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	url := app.mergePath(fmt.Sprintf("/v1/track/%d", id))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		app.errorResponse(w, r, res.StatusCode, http.StatusText(res.StatusCode))
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var payload envelope
	if err := json.Unmarshal(body, &payload); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"track": payload}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
