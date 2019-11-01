package types

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func load(name string) string {
	s, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func TestParseToken(t *testing.T) {
	s := load("fixtures/token.json")
	var token Token
	_ = json.Unmarshal([]byte(s), &token)
	expected := Token{Username: "toto", Version: "9.19.130", Token: "ABCD"}
	if !reflect.DeepEqual(token, expected) {
		t.Errorf("got %v != %v", token, expected)
	}
}

func TestParseTime(t *testing.T) {
	tm, _ := ParseTime("20191018000014")

	expected := time.Date(2019, time.October, 18, 0, 0, 14, 0, time.UTC)
	if !reflect.DeepEqual(expected, tm) {
		t.Errorf("got %#v != %#v", tm, expected)
	}
}

func TestParseTime2(t *testing.T) {
	tm, _ := ParseTime2("Nov 1, 2019 7:16:49 AM")
	tm2, _ := ParseTime2("Oct 31, 2019 10:00:00 PM")

	expected := time.Date(2019, time.November, 1, 7, 16, 49, 0, time.UTC)
	if !reflect.DeepEqual(expected, tm) {
		t.Errorf("got %#v != %#v", tm, expected)
	}
	expected = time.Date(2019, time.October, 31, 22, 00, 00, 0, time.UTC)
	if !reflect.DeepEqual(expected, tm2) {
		t.Errorf("got %#v != %#v", tm2, expected)
	}
}
