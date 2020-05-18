package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func createErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
	}
}

func parseJSONBody(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		log.Println(err)
		return errors.New("invalid request body")
	}

	return nil
}

func sendResponse(w http.ResponseWriter, status int, res interface{}) {
	//set the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the response
	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		log.Panic(err) //panic if can't write response
	}
}

func sendSuccessResponse(w http.ResponseWriter) {
	sendResponse(w, http.StatusOK, BasicResponse{Success: true})
}

func sendInternalErrorResponse(w http.ResponseWriter) {
	sendResponse(w, http.StatusInternalServerError, createErrorResponse("an internal error occurred"))
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
