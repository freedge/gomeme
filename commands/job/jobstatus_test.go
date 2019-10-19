package job

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

func TestJobStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		{
			"statuses" : [ {
			  "jobId" : "LUCCT1P:32zfh",
			  "folderId" : "LUCCT1P:",
			  "numberOfRuns" : 1,
			  "name" : "C1P-PRD-DBA-DBA",
			  "folder" : "C1P-PRD-DBA-DBA",
			  "type" : "Folder",
			  "status" : "Executing",
			  "held" : false,
			  "deleted" : false,
			  "startTime" : "20191018000014",
			  "endTime" : "",
			  "orderDate" : "191018",
			  "ctm" : "LUCCT1P",
			  "description" : "",
			  "host" : "",
			  "application" : "C1P-PRD",
			  "subApplication" : "C1P-PRD-DBA-DBA",
			  "outputURI" : "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/output",
			  "logURI" : "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/log"
			} ],
			"returned" : 1,
			"total" : 1
		  }		  
		`))
	}))
	defer ts.Close()

	var js JobsStatusCommand
	commands.Endpoint = ts.URL + "/api"
	qr, err := js.Run(nil)

	if err != nil {
		t.Error(err)
	}
	expected := types.JobsStatusReply{
		Statuses: []types.Status{
			types.Status{
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
	if !reflect.DeepEqual(qr, expected) {
		t.Errorf("got %#v != %#v", qr, expected)
	}
}
