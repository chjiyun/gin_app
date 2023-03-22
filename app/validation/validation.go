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
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	//至少要有2种不同类型的字符
	re := regexp.MustCompile(`^[A-Za-z0-9-=;':,.?"/~!@#$%^*()_+}{]{6,20}$`)
	if ok = re.MatchString(v); !ok {
		return false
	}
	reList := [3]*regexp.Regexp{
		regexp.MustCompile(`^[A-Za-z]+$`),
		regexp.MustCompile(`^[0-9]+$`),
		regexp.MustCompile(`^[-=;':,.?"/~!@#$%^*()_+}{]+$`),
	}
	for i := 0; i < len(reList); i++ {
		if ok := reList[i].MatchString(v); ok {
			return false
		}
	}
	return true
}
