package rctest

import (
	"github.com/maprost/restclient"
	"net/http"
)

func Status200() restclient.Result {
	return Status(http.StatusOK)
}

func Status204() restclient.Result {
	return Status(http.StatusNoContent)
}

func Status400() restclient.Result {
	return Status(http.StatusBadRequest)
}

func Status403() restclient.Result {
	return Status(http.StatusForbidden)
}

func Status404() restclient.Result {
	return Status(http.StatusNotFound)
}

func Status500() restclient.Result {
	return Status(http.StatusInternalServerError)
}

func Status501() restclient.Result {
	return Status(http.StatusNotImplemented)
}

func Status503() restclient.Result {
	return Status(http.StatusServiceUnavailable)
}

func Status(code int) restclient.Result {
	return restclient.Result{StatusCode: code}
}

func FailedResponse(code int, msg string) restclient.Result {
	return restclient.Result{StatusCode: code, ResponseError: msg}
}
