package deploy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freedge/gomeme/commands"
)

func TestDeployPut(t *testing.T) {
	chf := make(chan string, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, param, _ := r.FormFile("definitionsFile")
		chf <- param.Filename
		w.Write([]byte(`[ {
			"deploymentFile" : "foo.json",
			"successfulFoldersCount" : 1,
			"successfulSmartFoldersCount" : 0,
			"successfulSubFoldersCount" : 0,
			"successfulJobsCount" : 4,
			"successfulConnectionProfilesCount" : 0,
			"successfulDriversCount" : 0,
			"isDeployDescriptorValid" : false,
			"deployedFolders" : [ "bar" ]
		  } ]  
		`))
	}))
	defer ts.Close()

	var cmd put
	cmd.Filename = "fixtures/folder.json"
	cmd.Ctm = "workbench"
	defer func() { commands.Opts = commands.DefaultOpts{} }()
	commands.Opts.Endpoint = ts.URL + "/api"
	commands.Opts.Subject = "subject"

	err := cmd.Execute([]string{})

	if err != nil {
		t.Error(err)
	}
	if <-chf != "folder.json" {
		t.Errorf("not the proper file name in %#v", cmd.reply)
	}
	if cmd.reply[0].SuccessfulJobsCount != 4 {
		t.Errorf("not the proper job count in %#v", cmd.reply)
	}
}
