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

func Status401() restclient.Result {
	return Status(http.StatusUnauthorized)
}

func Status402() restclient.Result {
	return Status(http.StatusPaymentRequired)
}

func Status403() restclient.Result {
	return Status(http.StatusForbidden)
}

func Status404() restclient.Result {
	return Status(http.StatusNotFound)
}

func Status405() restclient.Result {
	return Status(http.StatusMethodNotAllowed)
}

func Status406() restclient.Result {
	return Status(http.StatusNotAcceptable)
}

func Status407() restclient.Result {
	return Status(http.StatusProxyAuthRequired)
}

func Status408() restclient.Result {
	return Status(http.StatusRequestTimeout)
}

func Status409() restclient.Result {
	return Status(http.StatusConflict)
}

func Status418() restclient.Result {
	return Status(http.StatusTeapot)
}

func Status500() restclient.Result {
	return Status(http.StatusInternalServerError)
}

func Status501() restclient.Result {
	return Status(http.StatusNotImplemented)
}

func Status502() restclient.Result {
	return Status(http.StatusBadGateway)
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
