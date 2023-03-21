package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ValidateFunc struct {
}

func (vf ValidateFunc) VerifyPhoneNumber(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	re := regexp.MustCompile(`^1[345789]\d{9}$`)
	return re.MatchString(v)
}

func (vf ValidateFunc) VerifyPassword(fl validator.FieldLevel) bool {

	return true
}
