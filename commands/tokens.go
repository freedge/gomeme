package commands

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/freedge/gomeme/types"
)

var (
	// Tokens to use to connect to the endpoint
	Tokens types.TokenFile = types.TokenFile{Endpoint: map[string]types.TokenFileToken{}}

	tokenFilePath string
)

const (
	// TokenFileName is the name of the token file under the user home directory
	tokenFileName = ".tokens"
)

// WriteTokensFile writes a file containing tokens
func WriteTokensFile() {
	var (
		bytes []byte
		err   error
	)

	if bytes, err = json.Marshal(Tokens); err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile(tokenFilePath, bytes, 0600); err != nil {
		panic(err)
	}
}

// ReadTokens sets tokens from tokenfile
func ReadTokens() {
	usr, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	tokenFilePath = filepath.Join(usr, tokenFileName)
	var bytes []byte
	if bytes, err = ioutil.ReadFile(tokenFilePath); err != nil {
		os.Remove(tokenFilePath)
		WriteTokensFile()
	} else if err = json.Unmarshal(bytes, &Tokens); err != nil {
		os.Remove(tokenFilePath)
		WriteTokensFile()
	}
	if Tokens.Endpoint == nil {
		Tokens.Endpoint = map[string]types.TokenFileToken{}
	}
}

// init sets token from tokenfile
func init() {
	ReadTokens()
}
