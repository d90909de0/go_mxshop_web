package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatorMobile(field validator.FieldLevel) bool {
	mobile := field.Field().String()
	matched, _ := regexp.Match(`^1\d{10}$`, []byte(mobile))
	return matched
}
