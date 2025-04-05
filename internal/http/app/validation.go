package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ValidationAPIError struct {
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type ValidationAPIErrors struct {
	Errors []ValidationAPIError `json:"errors"`
}

func (ve ValidationAPIErrors) Error() string {
	var output string
	for _, e := range ve.Errors {
		output += fmt.Sprintf("Field '%s' is invalid: %s\n", e.Source.Pointer, e.Detail)
	}

	return output
}

type GreenValidator struct {
	validator *validator.Validate
}

func NewGreenValidator() *GreenValidator {
  return &GreenValidator{
    validator: validator.New(),
  }
}

func (cv *GreenValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var validationError ValidationAPIErrors
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal validation error")
		}

		errors := make([]ValidationAPIError, 0, len(ve))

		for _, e := range ve {
			validationError := ValidationAPIError{
				Title:  "Validation Error",
				Detail: fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()),
			}
			validationError.Source.Pointer = strings.ToLower(fmt.Sprintf("/data/attributes/%s", e.Field()))

			errors = append(errors, validationError)
		}

		validationError.Errors = errors
		return validationError
	}

	return nil
}

