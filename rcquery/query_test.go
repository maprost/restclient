package rcquery_test

import (
	"github.com/maprost/assertion"
	"github.com/maprost/restclient/rcquery"
	"net/url"
	"testing"
)

func TestEmptyQuery(t *testing.T) {
	assert := assertion.New(t)

	query := rcquery.New().Get()
	assert.Equal(query, "")
}

func TestAllQueryTypes(t *testing.T) {
	assert := assertion.New(t)

	query := rcquery.New().
		Add("Blob", "Crop").
		Add("Limit", uint(23)).
		Add("Flag", true).
		Add("kilo", -12.1234567890).
		Add("pList", []int{1, 2, 3}).
		Add("crap", map[string]string{"not": "in"}).
		Get()
	assert.Equal(query, "?Blob=Crop&Limit=23&Flag=true&kilo=-12.123456789&pList=1&pList=2&pList=3")
}

func TestEscapeQuery(t *testing.T) {
	assert := assertion.New(t)

	query := rcquery.New().Add("Blob", "$ü!").Get()
	assert.Equal(query, "?Blob=%24%C3%BC%21")

	unescapedQuery, e := url.QueryUnescape(query)
	assert.Nil(e)
	assert.Equal(unescapedQuery, "?Blob=$ü!")
}

func TestNilTypes(t *testing.T) {
	assert := assertion.New(t)

	var str string = "nil"
	var i int = 0
	var j *int

	query := rcquery.New().
		Add("Query", &str).
		Add("PlainNil", nil).
		Add("FlagI", &i).
		Add("FlagJ", j).
		Get()
	assert.Equal(query, "?Query=nil&FlagI=0")
}
