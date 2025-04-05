package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIError struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func greenHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	var message interface{} = "Internal server error"
	if ve, ok := err.(ValidationAPIErrors); ok {
		code = http.StatusUnprocessableEntity
		c.JSON(code, ve)
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}
	c.Logger().Error(err)

	c.JSON(code, map[string]any{"errors": []APIError{
		{
			Code:    code,
			Message: message,
		},
	}})
}

