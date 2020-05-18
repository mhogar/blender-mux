package controllers

import (
	"log"
	"net/http"

	"blendermux/common"
	"blendermux/server/database"
	"blendermux/server/models"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type SessionController struct {
	database.UserCRUD
	database.SessionCRUD
}

//PostLogin handles Post requests to "/login"
func (con SessionController) PostLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) (int, interface{}) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//parse the body
	err := parseJSONBody(req.Body, &body)
	if err != nil {
		log.Println(common.ChainError("error parsing PostLogin request body", err))
		return http.StatusBadRequest, ErrorResponse{Success: false, Error: "invalid json body"}
	}

	//get the user
	user, _ := con.GetUserByUsername(body.Email)
	if user == nil {
		return http.StatusBadRequest, ErrorResponse{Success: false, Error: "invalid email and/or password"}
	}

	//validate the password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(body.Password))
	if err != nil {
		log.Println(common.ChainError("error comparing password hashes", err))
		return http.StatusBadRequest, ErrorResponse{Success: false, Error: "invalid email and/or password"}
	}

	session := models.CreateNewSession(user.ID)

	//generate a seesion cookie
	c := &http.Cookie{
		Name:  "session",
		Value: session.ID.String(),
	}
	http.SetCookie(w, c)

	//add session to db
	con.CreateSession(session)

	//return success
	return createSuccessResponse()
}

//PostLogout handles Post requests to "/logout"
func (con SessionController) PostLogout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) (int, interface{}) {
	//get the session id from the cookie
	sID, err := getUserSession(req)
	if err != nil {
		log.Println("error getting session from cookie")
		return http.StatusUnauthorized, createErrorResponse("invalid session")
	}

	//get the session
	session, err := con.GetSessionByID(sID)
	if err != nil {
		log.Println("no session found with id", sID.String())
		return http.StatusBadRequest, createErrorResponse("invalid session")
	}

	//check if session was found
	if session == nil {
		log.Println("no session found with id", sID.String())
		return http.StatusUnauthorized, createErrorResponse("invalid session")
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
	return createSuccessResponse()
}
