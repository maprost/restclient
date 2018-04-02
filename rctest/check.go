package rctest

import (
	"testing"

	"github.com/maprost/restclient"
	"github.com/maprost/should"
)

func CheckResult(t testing.TB, actual restclient.Result, expected restclient.Result) {
	should.BeNil(t, actual.Err)
	should.BeEqual(t, actual.StatusCode, expected.StatusCode)

	if len(expected.ResponseError) > 0 {
		should.BeEqual(t, actual.ResponseError, expected.ResponseError)
	}
}
