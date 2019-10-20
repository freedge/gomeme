package curl

import (
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
)

func TestCurl(t *testing.T) {
	commands.TheToken = "abc"
	commands.Endpoint = "https//toto/api"
	commands.Insecure = true

	c := curl{}
	s, _ := c.Run()

	if !reflect.DeepEqual(s, "curl -k -H 'Accept: application/json' -H 'Authorization: Bearer abc' https//toto/api") {
		t.Errorf("got %s", c.out)
	}
}
