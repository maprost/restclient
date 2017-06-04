package restclient_test

import (
	"encoding/json"
	"encoding/xml"
	"github.com/maprost/assertion"
	"github.com/maprost/restclient"
	"github.com/maprost/restclient/rctest"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func runServer(f http.HandlerFunc) (url string) {
	path := "/test"
	mux := http.NewServeMux()
	mux.HandleFunc(path, f)
	testServer := httptest.NewServer(mux)

	url = testServer.URL + path
	return
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

func TestSendBodyWithJsonPostRestClient_ok(t *testing.T) {
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

func TestSendBodyWithXMLPostRestClient_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg string
	}

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Body
		err := xml.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Body not readable", http.StatusBadRequest)
			return
		}

		if body.Msg != "Blob" {
			http.Error(w, "Msg is wrong", http.StatusBadRequest)
			return
		}

		x, _ := xml.Marshal(body)
		w.Header().Set("Content-Type", "application/xml")
		w.Write(x)
	})

	var res Body
	result := restclient.Post(url).AddXMLBody(Body{Msg: "Blob"}).SendAndGetXMLResponse(&res)
	rctest.AssertResult(assert, result, rctest.Status200())
	assert.Equal(res.Msg, "Blob")
}

func TestQueryParam_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit != 14 {
			http.Error(w, "Wrong limit", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Get(url).AddQueryParam("limit", 14).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestPointerQueryParam_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit != 14 {
			http.Error(w, "Wrong limit", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	limit := 14
	result := restclient.Get(url).AddQueryParam("limit", &limit).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestPointerInBody_ok(t *testing.T) {
	assert := assertion.New(t)

	type Body struct {
		Msg     *string
		Flag    *bool
		Limit   *int
		Changed *time.Time
	}

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body Body
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Body not readable", http.StatusBadRequest)
			return
		}

		if body.Msg != nil || body.Flag != nil || body.Limit != nil || body.Changed != nil {
			http.Error(w, "Nil value is not nil", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Get(url).AddJsonBody(Body{}).Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestAcceptedLanguageHeader_ok(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang != "da" {
			http.Error(w, "Wrong lang", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Get(url).AddHeader("Accept-Language", "da").Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}

func TestNoLogging(t *testing.T) {
	assert := assertion.New(t)

	url := runServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	result := restclient.Get(url).NoLogger().Send()
	rctest.AssertResult(assert, result, rctest.Status204())
}
