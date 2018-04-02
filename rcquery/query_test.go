package rcquery_test

import (
	"net/url"
	"testing"

	"github.com/maprost/restclient/rcquery"
	"github.com/maprost/should"
)

func TestEmptyQuery(t *testing.T) {
	query := rcquery.New().Get()
	should.BeEqual(t, query, "")
}

func TestAllQueryTypes(t *testing.T) {
	query := rcquery.New().
		Add("Blob", "Crop").
		Add("Limit", uint(23)).
		Add("Flag", true).
		Add("kilo", -12.1234567890).
		Add("pList", []int{1, 2, 3}).
		Add("crap", map[string]string{"not": "in"}).
		Get()
	should.BeEqual(t, query, "?Blob=Crop&Limit=23&Flag=true&kilo=-12.123456789&pList=1&pList=2&pList=3")
}

func TestEscapeQuery(t *testing.T) {
	query := rcquery.New().Add("Blob", "$ü!").Get()
	should.BeEqual(t, query, "?Blob=%24%C3%BC%21")

	unescapedQuery, e := url.QueryUnescape(query)
	should.BeNil(t, e)
	should.BeEqual(t, unescapedQuery, "?Blob=$ü!")
}

func TestNilTypes(t *testing.T) {
	var str string = "nil"
	var i int = 0
	var j *int

	query := rcquery.New().
		Add("Query", &str).
		Add("PlainNil", nil).
		Add("FlagI", &i).
		Add("FlagJ", j).
		Get()
	should.BeEqual(t, query, "?Query=nil&FlagI=0")
}
