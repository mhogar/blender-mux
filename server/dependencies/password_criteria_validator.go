package dependencies

import (
	"blendermux/server/models"
	"sync"
)

var passwordCriteriaValidatorOnce sync.Once
var passwordCriteriaValidator models.PasswordCriteriaValidator

// ResolvePasswordCriteriaValidator resolves the PasswordCriteriaValidator dependency.
// Only the first call to this function will create a new PasswordCriteriaValidator, after which it will be retrieved from the cache.
func ResolvePasswordCriteriaValidator() models.PasswordCriteriaValidator {
	passwordCriteriaValidatorOnce.Do(func() {
		passwordCriteriaValidator = models.StandardPasswordCriteriaValidator{}
	})
	return passwordCriteriaValidator
}
