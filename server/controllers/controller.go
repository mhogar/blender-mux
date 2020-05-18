package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// CreateRequestHandler fowards the request to handleFunc and writes its response using the response writer.
func CreateRequestHandler(handleFunc func(http.ResponseWriter, *http.Request, httprouter.Params) (status int, res interface{})) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		status, res := handleFunc(w, req, params)

		//catch any panics to prevent the server from crashing
		err := recover()
		if err != nil {
			log.Println(err)
			status, res = createInternalErrorResponse()
		}

		//set the header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		//write the response
		encoder := json.NewEncoder(w)
		err = encoder.Encode(res)
		if err != nil {
			log.Panic(err) //panic if can't write response
		}
	}
}
