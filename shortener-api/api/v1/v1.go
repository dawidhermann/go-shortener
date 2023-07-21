package v1

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &RequestError{Err: err, Status: status}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
