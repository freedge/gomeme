package qr

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func load(name string) string {
	s, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func TestQrCommand(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		[{"name":"PRD-SEV","ctm":"LUCCT4P","available":"9","max":10,"workloadPolicy":"N/A"}]	  
		`))
	}))
	defer ts.Close()

	var cmd QRCommand
	commands.Endpoint = ts.URL + "/api"
	qr, err := cmd.Run(nil)

	if err != nil {
		t.Error(err)
	}
	expected := []types.QR{types.QR{Available: "9", Ctm: "LUCCT4P", Max: 10, Name: "PRD-SEV"}}
	if !reflect.DeepEqual(qr, expected) {
		t.Errorf("got %#v != %#v", qr, expected)
	}
}
