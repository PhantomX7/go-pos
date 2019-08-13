package middleware

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"gopkg.in/go-playground/validator.v8"
)

// helper function

// error handling middleware binding translator
func validationErrorToText(e *validator.FieldError) string {
	e.Field = strcase.SnakeCase(e.Field)
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", e.Field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field, e.Param)
	case "unique":
		return fmt.Sprintf("%s already exist", e.Field)
	case "exist":
		return fmt.Sprintf("%s does not exist", e.Field)
	case "gte":
		return fmt.Sprintf("%s must equal to or greater than %s", e.Field, e.Param)
	}
	return fmt.Sprintf("%s is not valid", e.Field)
}
