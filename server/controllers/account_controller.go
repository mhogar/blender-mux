package controllers

import (
	"log"
	"net/http"

	"github.com/blendermux/server/dependencies"
	"github.com/blendermux/server/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AccountController struct {
	dependencies.UserCRUD
}

// PostAccount handles Post requests to "/account"
func (con *AccountController) PostAccount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		//TODO: handle error
		log.Println(err)
		sendResponse(w, errorResponse{false, "invalid json body"})
		return
	}

	//TODO: validate the input

	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		//TODO: handle error
		log.Println(err)
		sendResponse(w, errorResponse{false, "error generating password hash"})
		return
	}

	//save the user
	con.CreateUser(&models.User{
		primitive.NewObjectID(),
		body.Email,
		hash,
	})

	//return success
	sendResponse(w, basicResponse{true})
}
