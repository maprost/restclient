package restclient

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type ResponseItem struct {
	header http.Header
	body   []byte
	Result Result
}

func (r *ResponseItem) String() (output string) {
	if r.Result.Err != nil {
		return
	}

	output = string(r.body)
	return
}

func (r *ResponseItem) XML(output interface{}) {
	if r.Result.Err != nil {
		return
	}

	// set the output if there is something
	if r.Result.StatusCode == http.StatusOK {
		r.Result.Err = xml.Unmarshal(r.body, output)
	}

	return
}

func (r *ResponseItem) Json(output interface{}) {
	if r.Result.Err != nil {
		return
	}

	// set the output if there is something
	if r.Result.StatusCode == http.StatusOK {
		r.Result.Err = json.Unmarshal(r.body, output)
	}

	return
}

func (r *ResponseItem) Header(key string) (values []string, ok bool) {
	if r.Result.Err != nil {
		return
	}

	values, ok = r.header[key]
	return
}

func (r *ResponseItem) Error() error {
	return r.Result.Error()
}
