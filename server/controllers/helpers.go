package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type BasicResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func createErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
	}
}

func createSuccessResponse() (int, BasicResponse) {
	return http.StatusOK, BasicResponse{Success: true}
}

func createInternalErrorResponse() (int, ErrorResponse) {
	return http.StatusInternalServerError, createErrorResponse("an internal error occurred")
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
