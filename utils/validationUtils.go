package utils

import "github.com/go-playground/validator"

func IsValid[T any](t T) error {
	return validator.New().Struct(t)
}
