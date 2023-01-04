package restclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maprost/restclient/rcdep"
	"github.com/maprost/restclient/rcquery"
)

const jsonContentType = "application/json; charset=utf-8"
const xmlContentType = "application/xml; charset=utf-8"
const contentType = "Content-Type"

var DefaultLogger = log.New(os.Stdout, "", 0)

type noLogger struct{}

func (nl noLogger) Printf(format string, v ...interface{}) {}

type RestClient struct {
	log           rcdep.Logger
	requestPath   string
	requestMethod string
	requestBody   io.Reader
	header        map[string][]string
	query         rcquery.Query
	err           error
	httpClient    *http.Client
	basicAuthUser string
	basicAuthPW   string
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
		log:         noLogger{},
		requestPath: path,
		header:      make(map[string][]string),
	}
}

func (r *RestClient) AddLogger(logger rcdep.Logger) *RestClient {
	r.log = logger
	return r
}

func (r *RestClient) AddHttpClient(httpClient *http.Client) *RestClient {
	r.httpClient = httpClient
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
	if _, ok := r.header[key]; ok {
		// update
		r.header[key] = append(r.header[key], value)
	} else {
		// insert
		r.header[key] = []string{value}
	}
	return r
}

func (r *RestClient) AddBasicAuth(user string, pw string) *RestClient {
	r.basicAuthUser = user
	r.basicAuthPW = pw
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

// AddBody adds a []byte to the request body.
// Only usable in Post/Put requests.
func (r *RestClient) AddBody(input []byte, contentTypeValue string) *RestClient {
	// check for error
	if r.err != nil {
		return r
	}

	if r.err == nil {
		r.requestBody = bytes.NewReader(input)
		r.AddHeader(contentType, contentTypeValue)
	}

	return r
}

func (r *RestClient) Send() (result Result) {
	responseItem := r.send()
	result = responseItem.Result

	return
}

func (r *RestClient) SendAndGetResponseItem() ResponseItem {
	return r.send()
}

func (r *RestClient) SendAndGetResponse() (output string, result Result) {
	responseItem := r.send()

	output = responseItem.String()
	result = responseItem.Result
	return
}

func (r *RestClient) SendAndGetJsonResponse(output interface{}) (result Result) {
	responseItem := r.send()

	responseItem.Json(output)
	result = responseItem.Result
	return
}

func (r *RestClient) SendAndGetXMLResponse(output interface{}) (result Result) {
	responseItem := r.send()

	responseItem.XML(output)
	result = responseItem.Result
	return
}

func (r *RestClient) send() (responseItem ResponseItem) {
	if r.err != nil {
		responseItem.Result.Err = r.err
		return
	}

	// create request
	url := r.requestPath + r.query.Get()
	request, err := http.NewRequest(r.requestMethod, url, r.requestBody)
	if err != nil {
		responseItem.Result.Err = err
		return
	}

	// add header
	for key, values := range r.header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	if r.basicAuthUser != "" {
		request.SetBasicAuth(r.basicAuthUser, r.basicAuthPW)
	}

	// send request
	if r.httpClient == nil {
		r.httpClient = http.DefaultClient
	}
	start := time.Now()
	response, err := r.httpClient.Do(request)
	duration := time.Now().Sub(start)
	r.log.Printf("request [time: %v] %s:%s", duration, r.requestMethod, url)
	//r.log.Printf("request headers %v", request.Header)
	if err != nil {
		responseItem.Result.Err = err
		return
	}
	defer response.Body.Close()

	// show header
	r.log.Printf("response Status: %v", response.Status)
	r.log.Printf("response Headers: %v", response.Header)
	responseItem.header = response.Header

	// get body
	responseItem.body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		responseItem.Result.Err = err
		return
	}
	r.log.Printf("response Body: %v", string(responseItem.body))

	// set link + status
	responseItem.Result.Link = response.Request.RequestURI
	responseItem.Result.StatusCode = response.StatusCode

	// set responseError of failed response (status >= 400)
	if responseItem.Result.StatusCode >= http.StatusBadRequest {
		responseItem.Result.ResponseError = string(responseItem.body)
	}

	return
}
