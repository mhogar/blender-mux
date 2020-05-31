package controllers

import (
	"blendermux/common"
	"net/http"

	"github.com/google/uuid"
)

func getSessionFromRequest(req *http.Request) (uuid.UUID, error) {
	c, err := req.Cookie("session")
	if err != nil {
		return uuid.Nil, common.ChainError("error getting session cookie", err)
	}

	sID, err := uuid.Parse(c.Value)
	if err != nil {
		return uuid.Nil, common.ChainError("error parsing session id", err)
	}

	return sID, nil
}
