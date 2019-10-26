package job

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func TestJobTree(t *testing.T) {
	ch := make(chan types.JobsStatusReply, 6)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := json.Marshal(<-ch)
		w.Write(bytes)
	}))
	defer ts.Close()

	// we inject a few jobs
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "A"},
		{JobId: "B"},
		{JobId: "C"},
		{JobId: "D"},
		{JobId: "E"},
	}}

	// we reply with the dependency between the 2 jobs.
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "A"},
	}}
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "Z"},
		{JobId: "B"},
	}}
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "A"},
		{JobId: "B"},
		{JobId: "C"},
	}}
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "A"},
		{JobId: "Z"},
		{JobId: "C"},
		{JobId: "D"},
	}}
	ch <- types.JobsStatusReply{Statuses: []types.Status{
		{JobId: "E"},
	}}

	js := jobTreeCommand{}
	commands.Endpoint = ts.URL + "/api"
	err := js.Execute([]string{})

	if err != nil {
		t.Error(err)
	}

	expected := []treenode{
		{js.nodes["D"], 0},
		{js.nodes["A"], 1},
		{js.nodes["Z"], 1},
		{js.nodes["C"], 1},
		{js.nodes["B"], 2},
		{js.nodes["E"], 0},
	}
	if !reflect.DeepEqual(len(js.tree), 6) {
		t.Errorf("got %#v != %#v", js.tree, expected)
	}

}
