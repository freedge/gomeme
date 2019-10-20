[![Build Status](https://travis-ci.org/freedge/gomeme.svg?branch=master)](https://travis-ci.org/freedge/gomeme)
[![Build Status](https://dev.azure.com/freedge/freedge/_apis/build/status/freedge.gomeme?branchName=master)](https://dev.azure.com/freedge/freedge/_build/latest?definitionId=1&branchName=master)
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
gomeme login -user toto
gomeme qr -name PRD*

```

Usage on Windows:
```
$env:GOMEME_INSECURE="true"
...
```

Use the ```-dump``` option to output in a Go-like format. ```-json``` option outputs in json format.

## Commands

### login

Get a token for a user. Writes it into a .token file in the current directory.

The password must be provided either through the GOMEME_PASSWORD environment variable, or
through the terminal.

```
gomeme login -user toto
```

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

### qr.set

set a qr

```
gomeme qr.set -name DEV-FOO -ctm BARCT4T -max 5
```

### job.log

get the output or logs of a job id

```
gomeme job.log -jobid FOOCCT4P:5nq1c
gomeme job.log -jobid FOOCCT4P:5nq1c -output
```

### job.order

order a job.

```
gomeme job.order -ctm FOOCT4T -folder ABC-DEV-OND -hold -jobs dABC1
```

### job.action

hold/delete/undelete/confirm/setToOk a job

```
gomeme job.action -action delete -jobid FOOCT4T:3z553
```

### curl

Just outputs the curl command to run to target the API by hand

```
gomeme curl
```

## License

http://www.apache.org/licenses/LICENSE-2.0
