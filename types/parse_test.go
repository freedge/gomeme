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
