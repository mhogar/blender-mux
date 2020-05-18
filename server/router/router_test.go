package router_test

import (
	"blendermux/common"
	"blendermux/server/controllers"
	controllermocks "blendermux/server/controllers/mocks"
	"blendermux/server/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
	RequestHandler controllermocks.RequestHandler
	Router         *httprouter.Router
}

func (suite *RouterTestSuite) SetupTest() {
	suite.RequestHandler = controllermocks.RequestHandler{}
	suite.Router = router.CreateRouter(&suite.RequestHandler)
}

func (suite *RouterTestSuite) TestRouter_SendsInternalServerErrorOnPanic() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/user", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("PostUser", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
		panic("test panic handler")
	})

	//act
	res, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	var errRes controllers.ErrorResponse
	status := common.ParseResponse(&suite.Suite, res, &errRes)

	//assert
	suite.EqualValues(http.StatusInternalServerError, status)
	suite.False(errRes.Success)
	suite.Contains(errRes.Error, "an internal error occurred")
}

func (suite *RouterTestSuite) TestRouter_PostUserHandledByCorrectHandleFunction() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/user", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("PostUser", mock.Anything, mock.Anything, mock.Anything)

	//act
	_, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.RequestHandler.AssertCalled(suite.T(), "PostUser", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *RouterTestSuite) TestRouter_DeleteUserHandledByCorrectHandleFunction() {
	//arrange
	server := httptest.NewServer(suite.Router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/user/1", nil)
	suite.Require().NoError(err)

	suite.RequestHandler.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything)

	//act
	_, err = http.DefaultClient.Do(req)
	suite.Require().NoError(err)

	//assert
	suite.RequestHandler.AssertCalled(suite.T(), "DeleteUser", mock.Anything, mock.Anything, mock.MatchedBy(func(params httprouter.Params) bool {
		return params.ByName("id") == "1"
	}))
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, &RouterTestSuite{})
}
