package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomBinder struct
type CustomBinder struct {
	Binder echo.Binder
}

// Bind tries to bind request into interface, and if it does then validate it
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := cb.Binder.Bind(i, c); err != nil && !errors.Is(err, echo.ErrUnsupportedMediaType) {
		fmt.Println("Error binding request: ", err)
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}

// CustomValidator holds custom validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}
