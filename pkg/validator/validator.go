package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate runs struct tag validation and returns a field->message map,
// or nil if the struct is valid.
func Validate(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	for _, fe := range err.(validator.ValidationErrors) {
		errs[fe.Field()] = fmt.Sprintf("failed on '%s' validation", fe.Tag())
	}
	return errs
}
