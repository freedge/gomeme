package curl

import (
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func TestCurl(t *testing.T) {
	commands.Tokens.Endpoint["https//toto/api"] = types.TokenFileToken{Token: types.Token{Token: "abc"}}
	commands.Opts.Endpoint = "https//toto/api"

	c := curl{}
	_ = c.Execute([]string{})

	if !reflect.DeepEqual(c.out, "curl -H 'Accept: application/json' -H 'Authorization: Bearer abc' https//toto/api") {
		t.Errorf("got %#v", c.out)
	}
}
