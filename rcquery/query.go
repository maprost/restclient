package rcquery

import (
	"net/url"
	"reflect"
	"strconv"
)

type Query struct {
	query string
}

func New() *Query {
	return &Query{}
}

func (q *Query) Add(key string, value interface{}) *Query {
	// ignore nil values
	if value == nil {
		return q
	}

	var separator string
	if len(q.query) == 0 {
		separator = "?"
	} else {
		separator = "&"
	}

	val := reflect.Indirect(reflect.ValueOf(value))
	switch val.Kind() {
	case reflect.Bool:
		q.query += separator + key + "=" + strconv.FormatBool(val.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		q.query += separator + key + "=" + strconv.FormatInt(val.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		q.query += separator + key + "=" + strconv.FormatUint(val.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		q.query += separator + key + "=" + strconv.FormatFloat(val.Float(), 'f', -1, 64)

	case reflect.String:
		q.query += separator + key + "=" + url.QueryEscape(val.String())

	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			q.Add(key, val.Index(i).Interface())
		}
	default:
		// will ignored
	}
	return q
}

func (q *Query) Get() string {
	return q.query
}
