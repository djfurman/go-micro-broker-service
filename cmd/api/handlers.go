package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"), http.StatusBadGateway)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create JSON to send to Auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "  ")

	// Call the Auth microservice
	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer res.Body.Close()

	// Ensure we receive the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if res.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling Auth service"))
		return
	}

	// Create a variable into which we can read response.Body
	var jsonFromAuthService jsonResponse

	// decode the JSON from the authentication-service
	err = json.NewDecoder(res.Body).Decode(&jsonFromAuthService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromAuthService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// Now this is a valid login
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromAuthService.Data

	app.writeJSON(w, http.StatusOK, payload)
}
