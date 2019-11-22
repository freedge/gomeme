package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freedge/gomeme/commands"
)

const sampleDescription = `你好`

func TestAnnotation(t *testing.T) {
	ch := make(chan string, 2)
	defer func() { commands.Opts = commands.DefaultOpts{} }()

	fmt.Println(".")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		ch <- r.Header.Get("Annotation-Subject")
		ch <- r.Header.Get("Annotation-Description")
		w.Write([]byte{})
	}))
	defer ts.Close()
	commands.Opts.Debug = true
	commands.Opts.Endpoint = ts.URL + "/api"
	commands.Opts.Subject = "hello world"
	commands.Opts.Description = sampleDescription

	var output string
	err := Call("GET", "/toto", nil, map[string]string{}, &output)

	if err != nil {
		t.Error(err)
	}
	s := <-ch
	if s != "hello world" {
		t.Error("wrong subject")
	}
	s = <-ch
	if s != sampleDescription {
		t.Error("wrong description", s)
	}
}
