package restclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/mleuth/restclient/rcdep"
	"github.com/mleuth/restclient/rcquery"
	"github.com/mleuth/timeutil"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const jsonContentType = "application/json; charset=utf-8"
const xmlContentType = "application/xml; charset=utf-8"
const contentType = "Content-Type"

var defaultLogger = log.New(os.Stdout, "", 0)

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

func (r *RestClient) AddQueryParam(key string, value interface{}) *RestClient {
	r.query.Add(key, value)
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
		r.header[contentType] = jsonContentType
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
		r.header[contentType] = xmlContentType
	}

	return r
}

func (r *RestClient) Send() (statusCode int, responseError string, err error) {
	_, statusCode, responseError, err = r.send()
	return
}

func (r *RestClient) SendAndGetJsonResponse(output interface{}) (statusCode int, responseError string, err error) {
	body, statusCode, responseError, err := r.send()
	if err != nil {
		return
	}

	// set the output if there is something
	if statusCode == http.StatusOK {
		err = json.Unmarshal(body, output)
	}

	return
}

func (r *RestClient) SendAndGetXMLResponse(output interface{}) (statusCode int, responseError string, err error) {
	body, statusCode, responseError, err := r.send()
	if err != nil {
		return
	}

	// set the output if there is something
	if statusCode == http.StatusOK {
		err = xml.Unmarshal(body, output)
	}

	return
}

func (r *RestClient) send() (body []byte, statusCode int, responseError string, err error) {
	if r.err != nil {
		err = r.err
		return
	}

	// create request
	url := r.requestPath + r.query.Get()
	request, err := http.NewRequest(r.requestMethod, url, r.requestBody)
	if err != nil {
		return
	}

	// add header
	for key, value := range r.header {
		request.Header.Set(key, value)
	}

	// send request
	client := http.DefaultClient
	stopwatch := timeutil.NewStopwatch()
	response, err := client.Do(request)
	r.log.Printf("request [time: %v] %s:%s", stopwatch.String(), r.requestMethod, url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// show header
	r.log.Printf("response Status: %v", response.Status)
	r.log.Printf("response Headers: %v", response.Header)

	// get body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	r.log.Printf("response Body: %v", string(body))

	// set status
	statusCode = response.StatusCode

	// set responseError of failed response (status >= 400)
	if statusCode >= http.StatusBadRequest {
		responseError = string(body)
	}

	return
}
