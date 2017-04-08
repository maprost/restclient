package restclient_test

import (
	"encoding/json"
	"github.com/mleuth/assertion"
	"github.com/mleuth/restclient"
	"net/http"
	"net/http/httptest"
	"testing"
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

	statusCode, responseErr, err := restclient.Get(testServer.URL + "/test").Send()
	assert.Nil(err)
	assert.Equal(statusCode, http.StatusNoContent)
	assert.Equal(responseErr, "")
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

	var result Result
	statusCode, responseErr, err := restclient.Get(testServer.URL + "/test").SendAndGetJsonResponse(&result)
	assert.Nil(err)
	assert.Equal(statusCode, http.StatusOK)
	assert.Equal(responseErr, "")
	assert.Equal(result, Result{Msg: "Blob"})
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

	statusCode, responseErr, err := restclient.Get(testServer.URL + "/test").Send()
	assert.Nil(err)
	assert.Equal(statusCode, http.StatusBadRequest)
	assert.Equal(responseErr, "Blob is broken\n")
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

	statusCode, responseErr, err := restclient.Get(testServer.URL + "/test").AddJsonBody(Body{Msg: "Blob"}).Send()
	assert.Nil(err)
	assert.Equal(statusCode, http.StatusNoContent)
	assert.Equal(responseErr, "")
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

	statusCode, responseErr, err := restclient.Post(testServer.URL + "/test").AddJsonBody(Body{Msg: "Blob"}).Send()
	assert.Nil(err)
	assert.Equal(statusCode, http.StatusNoContent)
	assert.Equal(responseErr, "")
}
