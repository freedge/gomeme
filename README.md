[![Build Status](https://travis-ci.org/freedge/gomeme.svg?branch=master)](https://travis-ci.org/freedge/gomeme)
[![Go Report Card](https://goreportcard.com/badge/github.com/freedge/gomeme)](https://goreportcard.com/report/github.com/freedge/gomeme)
[![GoDoc](https://godoc.org/github.com/freedge/gomeme?status.svg)](https://godoc.org/github.com/freedge/gomeme) 


# gomeme



gomeme is a cli for Control-M, based on Control-M Automation API

https://docs.bmc.com/docs/automation-api/9181/services-784100995.html

and loosely inspired by govc.

## Usage

```
export GOMEME_ENDPOINT=https://.../automation-api
export GOMEME_PASSWORD=...
export GOMEME_INSECURE=true
gomeme login -username toto
gomeme qr -name PRD*

```

Usage on Windows:
```
$env:GOMEME_INSECURE="true"
...
```

Use the ```-dump``` option to output in a Go-like format. ```-json``` option should come next.

## Commands

### qr

list qrs

```
gomeme qr -name PRD-*
```

### lj

list jobs (default limit is)

```
gomeme lj -application TOTO-PRD -status Executing -limit 30
```


## License

http://www.apache.org/licenses/LICENSE-2.0
