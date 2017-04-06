package rcquery

import (
	"fmt"
	"net/url"
	"reflect"
)

type Query struct {
	query string
}

func New() *Query {
	return &Query{}
}

func (q *Query) Add(key string, value interface{}) *Query {
	var separator string
	if len(q.query) == 0 {
		separator = "?"
	} else {
		separator = "&"
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		q.query += separator + key + "=" + fmt.Sprint(value)

	case reflect.String:
		q.query += separator + key + "=" + url.QueryEscape(fmt.Sprint(value))

	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			colVal := val.Index(i).Interface()
			q.Add(key, colVal)
		}
	default:
		// will ignored
	}
	return q
}

func (q *Query) Get() string {
	return q.query
}
