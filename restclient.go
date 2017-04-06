package restclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/mleuth/restclient/rcdep"
	"github.com/mleuth/restclient/rcquery"
	"github.com/mleuth/timeutil"
	"io"
	"io/ioutil"
	"net/http"
)

const jsonContentType = "application/json; charset=utf-8"
const xmlContentType = "application/xml; charset=utf-8"
const contentType = "Content-Type"
const RestClientError = -1

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
	return &RestClient{log: defaultLogger, requestPath: path, requestMethod: http.MethodGet}
}

func Post(path string) *RestClient {
	return &RestClient{log: defaultLogger, requestPath: path, requestMethod: http.MethodPost}
}

func Put(path string) *RestClient {
	return &RestClient{log: defaultLogger, requestPath: path, requestMethod: http.MethodPut}
}

func Delete(path string) *RestClient {
	return &RestClient{log: defaultLogger, requestPath: path, requestMethod: http.MethodDelete}
}

func (r *RestClient) AddLogger(logger rcdep.Logger) *RestClient {
	r.log = logger
	return r
}

func (r *RestClient) AddQueryParam(key string, value interface{}) *RestClient {
	r.query.Add(key, value)
	return r
}

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

func (r *RestClient) Send() (statusCode int, err error) {
	_, statusCode, err = r.send()
	return
}

func (r *RestClient) SendAndGetJsonResponse(output interface{}) (statusCode int, err error) {
	body, statusCode, err := r.send()
	if statusCode == RestClientError {
		return statusCode, err
	}

	// set the output if there is something
	if statusCode == http.StatusOK {
		err = json.Unmarshal(body, output)
		if err != nil {
			return RestClientError, err
		}
	}

	return
}

func (r *RestClient) SendAndGetXMLResponse(output interface{}) (statusCode int, err error) {
	body, statusCode, err := r.send()
	if statusCode == RestClientError {
		return statusCode, err
	}

	// set the output if there is something
	if statusCode == http.StatusOK {
		err = xml.Unmarshal(body, output)
		if err != nil {
			return RestClientError, err
		}
	}

	return
}

func (r *RestClient) send() (body []byte, statusCode int, err error) {
	if r.err != nil {
		statusCode = RestClientError
		err = r.err
		return
	}

	url := r.requestPath + r.query.Get()
	request, err := http.NewRequest(r.requestMethod, url, r.requestBody)
	if err != nil {
		statusCode = RestClientError
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
	r.log.Printf("request [time: "+stopwatch.String()+"] "+r.requestMethod, ":", url)

	if err != nil {
		statusCode = RestClientError
		return
	}
	defer response.Body.Close()

	// show header
	r.log.Printf("response Status: %v", response.Status)
	r.log.Printf("response Headers: %v", response.Header)

	// get body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		statusCode = RestClientError
		return
	}

	// log body
	r.log.Printf("response Body: %v", string(body))

	// set status
	statusCode = response.Status

	// set error of failed response (status >= 400)
	if response.Status >= http.StatusBadRequest {
		err = errors.New(string(body))
	}

	return
}
