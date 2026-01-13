package validations

import (
	"regexp"
	"unicode"
)

// Regex for a simple email validation (lowercase local and domain)
// Use a single backslash before the dot in the raw string so the
// regex engine receives "\." to match a literal period.
var Regex_correo = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidatePassword(s string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	// Require a minimum length (8) but no strict upper limit here.
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
			//case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}
