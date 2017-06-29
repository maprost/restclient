package rctest

import (
	"github.com/maprost/assertion"
	"testing"
)

func TestStatus(t *testing.T) {
	assert := assertion.New(t)

	assert.Equal(Status200().StatusCode, 200)
	assert.Equal(Status204().StatusCode, 204)
	assert.Equal(Status400().StatusCode, 400)
	assert.Equal(Status403().StatusCode, 403)
	assert.Equal(Status404().StatusCode, 404)
	assert.Equal(Status500().StatusCode, 500)
	assert.Equal(Status501().StatusCode, 501)
	assert.Equal(Status503().StatusCode, 503)
}
