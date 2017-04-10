package restclient

import (
	"errors"
	"strconv"
)

type Result struct {
	StatusCode    int
	ResponseError string
	Err           error
}

func (r Result) Error() error {
	if r.Err != nil {
		return r.Err
	}
	if r.StatusCode >= 400 {
		return errors.New("[" + strconv.Itoa(r.StatusCode) + "]" + r.ResponseError)
	}
	return nil
}
