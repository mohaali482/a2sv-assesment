package infrastructure

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(validate *validator.Validate, i interface{}) error {
	RegisterTagNameFunc(validate)
	return validate.Struct(i)
}

func RegisterTagNameFunc(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ReturnErrorResponse(err error) []ErrorResponse {
	var validationErrors []ErrorResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ErrorResponse{
				Field:   err.Field(),
				Message: SetValidationResult(err.Tag()),
			})
		}
	}
	return validationErrors
}

func SetValidationResult(tag string) string {
	switch tag {
	case "required":
		tag = "This field is required"
	case "email":
		tag = "This field is not a valid email"
	case "min":
		tag = "This field is too short"
	case "max":
		tag = "This field is too long"
	case "eqfield":
		tag = "This field is not equal to the other field"
	case "mongodb":
		tag = "This field is not a valid ObjectID"
	}
	return tag
}
