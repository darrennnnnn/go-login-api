package validation

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", Password)
	}
}

// Password validations
var (
	upper   = regexp.MustCompile(`[A-Z]`)
	lower   = regexp.MustCompile(`[a-z]`)
	digit   = regexp.MustCompile(`[0-9]`)
	special = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

func Password(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	return len(password) >= 8 &&
		upper.MatchString(password) &&
		lower.MatchString(password) &&
		digit.MatchString(password) &&
		special.MatchString(password)
}