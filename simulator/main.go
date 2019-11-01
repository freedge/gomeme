package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
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
		},
		{
			"jobId" : "LUCCT1P:32zfh",
			"folderId" : "LUCCT1P:",
			"numberOfRuns" : 1,
			"name" : "C1P-PRD-DBA-DBA",
			"folder" : "C1P-PRD-DBA-DBA",
			"type" : "Folder",
			"status" : "Ended OK",
			"held" : false,
			"deleted" : false,
			"startTime" : "20191018000014",
			"endTime" : "",
			"orderDate" : "191018",
			"ctm" : "LUCCT1P",
			"description" : "",
			"host" : "mysuperhost",
			"application" : "C1P-PRD",
			"subApplication" : "C1P-PRD-DBA-DBA",
			"outputURI" : "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/output",
			"logURI" : "https://bla:8443/automation-api/run/job/LUCCT1P:32zfh/log"
		  },
		  {
			"jobId" : "LUCCT1P:32zfh",
			"folderId" : "LUCCT1P:",
			"numberOfRuns" : 1,
			"name" : "C1P-PRD-DBA-DBA",
			"folder" : "C1P-PRD-DBA-DBA",
			"type" : "Folder",
			"status" : "Ended Not OK",
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
		  }
		],
		"returned" : 3,
		"total" : 5
	  }		  
	`))
	}))
	defer ts.Close()

	fmt.Println(ts.URL)
	c := make(chan os.Signal, 1)
	for {
		select {
		case s := <-c:
			if s == syscall.SIGTERM {
				return
			}
		}

	}
}
