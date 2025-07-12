package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on its 'validate' tags.
func ValidateStruct(s interface{}) []string {
	if err := validate.Struct(s); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed on the '%s' tag", err.Field(), err.Tag()))
		}
		return errors
	}
	return nil
}
