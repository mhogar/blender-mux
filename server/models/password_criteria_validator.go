package models

// PasswordCriteriaValidator is an interface for validating a password against criteria.
type PasswordCriteriaValidator interface {
	// ValidatePasswordCriteria validates the password meets the minimum complexity criteria.
	ValidatePasswordCriteria(password string) error
}

// StandardPasswordCriteriaValidator is an implementation of PasswordCriteriaValidator that uses standard criteria.
type StandardPasswordCriteriaValidator struct{}

// ValidatePasswordCriteria validates the password meets the standard minimum complexity criteria.
func (StandardPasswordCriteriaValidator) ValidatePasswordCriteria(password string) error {
	//TODO: implement
	return nil
}
