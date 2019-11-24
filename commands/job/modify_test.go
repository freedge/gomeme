package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

const (
	initialJob = `{"dToto":{"Type":"Job:Script","DaysKeepActive":"1","SubApplication":"FOO","Priority":"sa","FileName":"bar-job.sh","FilePath":"/controlm","Description":"1234","RunAs":"controlm","TimeZone":"GMT","Application":"bar","Arguments":["X","XX","XXX"],"DocumentationUrl":{"Url":"https://example.com"},"RerunLimit":{"Times":"1"},"When":{"ToTime":">","FromTime":"0100"},"INIT":{"Type":"Resource:Semaphore","Quantity":"1"},"IfBase:Folder:CompletionStatus_0":{"Type":"If:CompletionStatus","CompletionStatus":"1","Action:SetToNotOK_0":{"Type":"Action:SetToNotOK"}}}}`
	finalJob   = `{"dToto":{"Type":"Job:Script","DaysKeepActive":"1","SubApplication":"FOO","Priority":"sa","FileName":"bar-job.sh","FilePath":"/controlm","Description":"1234","RunAs":"controlm","TimeZone":"GMT","Application":"bar","Arguments":["a","b","c"],"DocumentationUrl":{"Url":"https://example.com"},"RerunLimit":{"Times":"1"},"When":{"ToTime":">","FromTime":"0100"},"INIT":{"Type":"Resource:Semaphore","Quantity":"1"},"IfBase:Folder:CompletionStatus_0":{"Type":"If:CompletionStatus","CompletionStatus":"1","Action:SetToNotOK_0":{"Type":"Action:SetToNotOK"}}}}`
)

func TestModify(t *testing.T) {
	ch := make(chan []byte, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 0 {
			r.ParseMultipartForm(0)
			file, _, _ := r.FormFile("jobDefinitionsFile")
			bodybytes, _ := ioutil.ReadAll(file)
			ch <- bodybytes
			w.Write([]byte("modified"))
		} else {
			w.Write([]byte(initialJob))
		}
	}))
	defer ts.Close()

	js := modify{Jobid: "abc", Name: "dToto"}
	commands.Opts.Endpoint = ts.URL + "/api"
	defer func() { commands.Opts = commands.DefaultOpts{} }()
	commands.Opts.Subject = "1234"
	err := js.Execute([]string{"a", "b", "c"})

	if err != nil {
		t.Error(err)
	}
	bytes := <-ch
	var job types.JobGetReply
	json.Unmarshal(bytes, &job)
	fmt.Println(string(bytes))
	if !reflect.DeepEqual(job["dToto"].Arguments, []string{"a", "b", "c"}) {
		t.Errorf("invalid arguments received")
	}
	if string(bytes) != finalJob {
		t.Errorf("%s != %s", finalJob, string(bytes))
	}
}
