package dependencies

import (
	"blendermux/server/controllers"
	"sync"
)

var passwordHasherOnce sync.Once
var passwordHasher controllers.PasswordHasher

// ResolvePasswordHasher resolves the PasswordHasher dependency.
// Only the first call to this function will create a new PasswordHasher, after which it will be retrieved from the cache.
func ResolvePasswordHasher() controllers.PasswordHasher {
	passwordHasherOnce.Do(func() {
		passwordHasher = controllers.BCryptPasswordHasher{}
	})
	return passwordHasher
}
