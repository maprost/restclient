package restclient_test

import (
	"encoding/json"
	"github.com/maprost/assertion"
	"github.com/maprost/restclient"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/maprost/restclient/rctest"
)

func Test204GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	getRequest := false
	// setup function to test
	test := func(w http.ResponseWriter, r *http.Request) {
		getRequest = r.Method == http.MethodGet
		w.WriteHeader(http.StatusNoContent)
	}

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/test", test)
	testServer := httptest.NewServer(mux)

	result := restclient.Get(testServer.URL + "/test").Send()
	rctest.AssertResult(assert, result, rctest.Status204())
	assert.True(getRequest)
}

func Test200GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Result struct {
		Msg string
	}

	// setup function to test
	test := func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(Result{Msg: "Blob"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/test", test)
	testServer := httptest.NewServer(mux)

	var r Result
	result := restclient.Get(testServer.URL + "/test").SendAndGetJsonResponse(&r)
	rctest.AssertResult(assert, result, rctest.Status200())
	assert.Equal(r, Result{Msg: "Blob"})
}

func Test404GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	// setup function to test
	test := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Blob is broken", http.StatusBadRequest)
	}

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/test", test)
	testServer := httptest.NewServer(mux)

	result := restclient.Get(testServer.URL + "/test").Send()
	rctest.AssertResult(assert, result, rctest.FailedResponse(400, "Blob is broken\n"))
}

func TestSendBodyWithGetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg string
	}

	// setup function to test
	test := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Body
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Body not readable", http.StatusBadRequest)
			return
		}

		if body.Msg != "Blob" {
			http.Error(w, "Msg is wrong", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/test", test)
	testServer := httptest.NewServer(mux)

	result := restclient.Get(testServer.URL + "/test").AddJsonBody(Body{Msg: "Blob"}).Send()
	assert.Nil(result.Err)
	assert.Equal(result.StatusCode, http.StatusNoContent)
	assert.Equal(result.ResponseError, "")
}

func TestSendBodyWithPostRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg string
	}

	// setup function to test
	test := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Body
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Body not readable", http.StatusBadRequest)
			return
		}

		if body.Msg != "Blob" {
			http.Error(w, "Msg is wrong", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/test", test)
	testServer := httptest.NewServer(mux)

	result := restclient.Post(testServer.URL + "/test").AddJsonBody(Body{Msg: "Blob"}).Send()
	assert.Nil(result.Err)
	assert.Equal(result.StatusCode, http.StatusNoContent)
	assert.Equal(result.ResponseError, "")
}
