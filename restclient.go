package restclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/maprost/restclient/rcdep"
	"github.com/maprost/restclient/rcquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const jsonContentType = "application/json; charset=utf-8"
const xmlContentType = "application/xml; charset=utf-8"
const contentType = "Content-Type"

var defaultLogger = log.New(os.Stdout, "", 0)

type noLogger struct{}

func (nl noLogger) Printf(format string, v ...interface{}) {}

type RestClient struct {
	log           rcdep.Logger
	requestPath   string
	requestMethod string
	requestBody   io.Reader
	header        map[string]string
	query         rcquery.Query
	err           error
}

func Get(path string) *RestClient {
	rc := newRC(path)
	rc.requestMethod = http.MethodGet
	return rc
}

func Post(path string) *RestClient {
	rc := newRC(path)
	rc.requestMethod = http.MethodPost
	return rc
}

func Put(path string) *RestClient {
	rc := newRC(path)
	rc.requestMethod = http.MethodPut
	return rc
}

func Delete(path string) *RestClient {
	rc := newRC(path)
	rc.requestMethod = http.MethodDelete
	return rc
}

func newRC(path string) *RestClient {
	return &RestClient{
		log:         defaultLogger,
		requestPath: path,
		header:      map[string]string{},
	}
}

func (r *RestClient) AddLogger(logger rcdep.Logger) *RestClient {
	r.log = logger
	return r
}

func (r *RestClient) NoLogger() *RestClient {
	r.log = noLogger{}
	return r
}

func (r *RestClient) AddQueryParam(key string, value interface{}) *RestClient {
	r.query.Add(key, value)
	return r
}

func (r *RestClient) AddHeader(key string, value string) *RestClient {
	r.header[key] = value
	return r
}

// AddJsonBody adds a struct as json to the request body.
// Only usable in Post/Put requests.
func (r *RestClient) AddJsonBody(input interface{}) *RestClient {
	// check for error
	if r.err != nil {
		return r
	}

	js := new(bytes.Buffer)
	r.err = json.NewEncoder(js).Encode(input)

	if r.err == nil {
		r.requestBody = js
		r.AddHeader(contentType, jsonContentType)
	}

	return r
}

// AddXMLBody adds a struct as xml to the request body.
// Only usable in Post/Put requests.
func (r *RestClient) AddXMLBody(input interface{}) *RestClient {
	// check for error
	if r.err != nil {
		return r
	}

	x := new(bytes.Buffer)
	r.err = xml.NewEncoder(x).Encode(input)

	if r.err == nil {
		r.requestBody = x
		r.AddHeader(contentType, xmlContentType)
	}

	return r
}

func (r *RestClient) Send() (result Result) {
	_, result = r.send()
	return
}

func (r *RestClient) SendAndGetJsonResponse(output interface{}) (result Result) {
	body, result := r.send()
	if result.Err != nil {
		return
	}

	// set the output if there is something
	if result.StatusCode == http.StatusOK {
		result.Err = json.Unmarshal(body, output)
	}

	return
}

func (r *RestClient) SendAndGetXMLResponse(output interface{}) (result Result) {
	body, result := r.send()
	if result.Err != nil {
		return
	}

	// set the output if there is something
	if result.StatusCode == http.StatusOK {
		result.Err = xml.Unmarshal(body, output)
	}

	return
}

func (r *RestClient) send() (body []byte, result Result) {
	if r.err != nil {
		result.Err = r.err
		return
	}

	// create request
	url := r.requestPath + r.query.Get()
	request, err := http.NewRequest(r.requestMethod, url, r.requestBody)
	if err != nil {
		result.Err = err
		return
	}

	// add header
	for key, value := range r.header {
		request.Header.Set(key, value)
	}

	// send request
	client := http.DefaultClient
	start := time.Now()
	response, err := client.Do(request)
	duration := time.Now().Sub(start)
	r.log.Printf("request [time: %v] %s:%s", duration, r.requestMethod, url)
	if err != nil {
		result.Err = err
		return
	}
	defer response.Body.Close()

	// show header
	r.log.Printf("response Status: %v", response.Status)
	r.log.Printf("response Headers: %v", response.Header)

	// get body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		result.Err = err
		return
	}
	r.log.Printf("response Body: %v", string(body))

	// set status
	result.StatusCode = response.StatusCode

	// set responseError of failed response (status >= 400)
	if result.StatusCode >= http.StatusBadRequest {
		result.ResponseError = string(body)
	}

	return
}
