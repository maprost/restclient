package rctest

import (
	"github.com/mleuth/assertion"
	"github.com/mleuth/restclient"
)

func CheckResult(t assertion.TestEnvironment, actual restclient.Result, expected restclient.Result) {
	assert := assertion.New(t)
	AssertResult(assert, actual, expected)
}

func AssertResult(assert assertion.Assert, actual restclient.Result, expected restclient.Result) {
	assert.Nil(actual.Err)
	assert.Equal(actual.StatusCode, expected.StatusCode)
	if len(expected.ResponseError) > 0 {
		assert.Equal(actual.ResponseError, expected.ResponseError)
	}
}
