package rctest

import (
	"testing"

	"github.com/maprost/should"
)

func TestStatus(t *testing.T) {
	should.BeEqual(t, Status200().StatusCode, 200)
	should.BeEqual(t, Status204().StatusCode, 204)

	should.BeEqual(t, Status400().StatusCode, 400)
	should.BeEqual(t, Status401().StatusCode, 401)
	should.BeEqual(t, Status402().StatusCode, 402)
	should.BeEqual(t, Status403().StatusCode, 403)
	should.BeEqual(t, Status404().StatusCode, 404)
	should.BeEqual(t, Status405().StatusCode, 405)
	should.BeEqual(t, Status406().StatusCode, 406)
	should.BeEqual(t, Status407().StatusCode, 407)
	should.BeEqual(t, Status408().StatusCode, 408)
	should.BeEqual(t, Status409().StatusCode, 409)
	should.BeEqual(t, Status418().StatusCode, 418)

	should.BeEqual(t, Status500().StatusCode, 500)
	should.BeEqual(t, Status501().StatusCode, 501)
	should.BeEqual(t, Status502().StatusCode, 502)
	should.BeEqual(t, Status503().StatusCode, 503)
}
