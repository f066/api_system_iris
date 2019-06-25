package utils

import (
	"gopkg.in/go-playground/validator.v9"
)

var (
	Validate *validator.Validate
)

func init() {
	Validate = validator.New()
}

func errorData(errs ...error) string {
	var s string
	for _, err := range errs {
		if err != nil {
			s += err.Error() + "\n"
		}
	}
	return s
}