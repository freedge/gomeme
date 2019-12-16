package curl

import (
	"reflect"
	"testing"

	"github.com/freedge/gomeme/commands"
	"github.com/freedge/gomeme/types"
)

func TestInvokeRestMethod(t *testing.T) {
	commands.Tokens.Endpoint["https//toto/api"] = types.TokenFileToken{Token: types.Token{Token: "abc"}}
	commands.Opts.Endpoint = "https//toto/api"

	c := invokeRestMethod{}
	_ = c.Execute([]string{})

	if !reflect.DeepEqual(c.out, "Invoke-RestMethod -Headers @{Accept='application/json'; Authorization='Bearer abc'}  -Uri https//toto/api") {
		t.Errorf("got %#v", c.out)
	}
}
