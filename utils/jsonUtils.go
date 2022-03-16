package utils

import (
	"encoding/json"
	"io"
)

func ToJson[T any](t T, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(t)
}

func FromJson[T any](t T, reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&t)
}
