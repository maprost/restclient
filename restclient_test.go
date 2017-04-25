package restclient_test

import (
	"encoding/json"
	"github.com/maprost/assertion"
	"github.com/maprost/restclient"
	"github.com/maprost/restclient/rctest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func runServer(f http.HandlerFunc) (url string) {
	path := "/test"
	mux := http.NewServeMux()
	mux.HandleFunc(path, f)
	testServer := httptest.NewServer(mux)

	return testServer.URL + path
}

func Test204GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "No get method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Get(url).Send()
	assert.Nil(result.Error())
}

func Test200GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Result struct {
		Msg string
	}

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(Result{Msg: "Blob"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})

	var r Result
	result := restclient.Get(url).SendAndGetJsonResponse(&r)
	rctest.AssertResult(assert, result, rctest.Status200())
	assert.Equal(r, Result{Msg: "Blob"})
}

func Test404GetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Blob is broken", http.StatusBadRequest)
	})

	result := restclient.Get(url).Send()
	rctest.AssertResult(assert, result, rctest.FailedResponse(400, "Blob is broken\n"))
}

func TestSendBodyWithGetRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg string
	}

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
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
	})

	result := restclient.Get(url).AddJsonBody(Body{Msg: "Blob"}).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestSendBodyWithPostRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg string
	}

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
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
	})

	result := restclient.Post(url).AddJsonBody(Body{Msg: "Blob"}).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestPutRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "No put method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Put(url).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestDeleteRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "No delete method", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Delete(url).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}
