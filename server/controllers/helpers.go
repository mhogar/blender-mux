package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type basicResponse struct {
	Success bool `json:"success"`
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type dataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ParseJSONBody parses the body of req and stores the data in v
func parseJSONBody(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		log.Println(err)
		return errors.New("invalid request body")
	}

	return nil
}

func sendResponse(w http.ResponseWriter, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Println(err)
		return errors.New("failed to create response")
	}

	return nil
}

func getUserSession(req *http.Request) (uuid.UUID, error) {
	c, err := req.Cookie("session")
	if err != nil {
		return uuid.Nil, err
	}

	sID, err := uuid.Parse(c.Value)
	if err != nil {
		return uuid.Nil, err
	}

	return sID, nil
}