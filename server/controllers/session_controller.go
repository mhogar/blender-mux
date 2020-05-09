package controllers

import (
	"log"
	"net/http"

	"github.com/blendermux/server/database"
	"github.com/blendermux/server/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type SessionController struct {
	database.UserCRUD
	database.SessionCRUD
}

//PostLogin handles Post requests to "/login"
func (con SessionController) PostLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	//get the user
	user, _ := con.GetUserByEmail(body.Email)
	if user == nil {
		sendResponse(w, errorResponse{false, "Invalid email and/or password"})
		return
	}

	//validate the password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(body.Password))
	if err != nil {
		sendResponse(w, errorResponse{false, "Invalid email and/or password"})
		return
	}

	//generate a seesion cookie
	sID := uuid.New()
	c := &http.Cookie{
		Name:  "session",
		Value: sID.String(),
	}
	http.SetCookie(w, c)

	//add session to db
	con.CreateSession(&models.Session{
		sID,
		user.ID,
	})

	//return success
	sendResponse(w, basicResponse{true})
}

//PostLogout handles Post requests to "/logout"
func (con SessionController) PostLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//get the session
	sID, err := getUserSession(req)
	if err != nil {
		sendResponse(w, errorResponse{false, "user session is invalid"})
		return
	}

	//validate sID
	session, _ := con.GetSessionByID(sID)
	if session == nil {
		log.Println("no session found in db with id", sID.String())
		sendResponse(w, errorResponse{false, "user session is invalid"})
		return
	}

	//delete session from db
	con.DeleteSession(session)

	//remove the cookie
	c := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	//return success
	sendResponse(w, basicResponse{true})
}
