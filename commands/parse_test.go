package commands

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

func TestParseQr(t *testing.T) {
	s := load("fixtures/qr.json")
	var qr []QR
	_ = json.Unmarshal([]byte(s), &qr)
	expected := []QR{QR{Available: "9", Ctm: "LUCCT4P", Max: 10, Name: "PRD-SEV"}}
	if !reflect.DeepEqual(qr, expected) {
		t.Errorf("got %v != %v", qr, expected)
	}
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

func TestParseJobsStatuses(t *testing.T) {
	s := load("fixtures/status.json")
	var reply JobsStatusReply
	_ = json.Unmarshal([]byte(s), &reply)
	expected := JobsStatusReply{
		Statuses: []Status{
			Status{
				JobId:          "LUCCT1P:32zfh",
				FolderId:       "LUCCT1P:",
				NumberOfRuns:   1,
				Name:           "C1P-PRD-DBA-DBA",
				Folder:         "C1P-PRD-DBA-DBA",
				Type:           "Folder",
				Status:         "Executing",
				Held:           false,
				Deleted:        false,
				StartTime:      "20191018000014",
				EndTime:        "",
				OrderDate:      "191018",
				Ctm:            "LUCCT1P",
				Description:    "",
				Host:           "",
				Application:    "C1P-PRD",
				SubApplication: "C1P-PRD-DBA-DBA",
				OutputURI:      "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/output",
				LogURI:         "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/log",
			},
		}, Returned: 1, Total: 1}
	if !reflect.DeepEqual(reply, expected) {
		t.Errorf("got %#v != %#v", reply, expected)
	}
}

func TestParseTime(t *testing.T) {
	tm, _ := ParseTime("20191018000014")

	expected := time.Date(2019, time.October, 18, 0, 0, 14, 0, time.UTC)
	if !reflect.DeepEqual(expected, tm) {
		t.Errorf("got %#v != %#v", tm, expected)
	}
}
