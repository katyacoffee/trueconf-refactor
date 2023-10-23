package model

import (
	"errors"

	"github.com/go-chi/render"
)

var (
	UserNotFound = errors.New("user_not_found")
)

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
