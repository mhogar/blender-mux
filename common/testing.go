package common

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/suite"
)

// ParseResponse is a helper method to use in testing to parse an http response into its status code and body.
func ParseResponse(suite *suite.Suite, res *http.Response, body interface{}) (status int) {
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(body)
	suite.Require().NoError(err)

	return res.StatusCode
}
