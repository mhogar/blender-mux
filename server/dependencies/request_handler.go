package dependencies

import (
	"blendermux/server/controllers"
	"sync"
)

var requestHandlerOnce sync.Once
var requestHandler controllers.RequestHandler

// ResolveRequestHandler resolves the RouteHandler dependency.
// Only the first call to this function will create a new RouteHandler, after which it will be retrieved from the cache.
func ResolveRequestHandler() controllers.RequestHandler {
	requestHandlerOnce.Do(func() {
		requestHandler = &controllers.RequestHandle{
			UserController: controllers.UserController{
				UserCRUD:                  ResolveDatabase(),
				PasswordHasher:            ResolvePasswordHasher(),
				PasswordCriteriaValidator: ResolvePasswordCriteriaValidator(),
			},
		}
	})
	return requestHandler
}
