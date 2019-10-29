[![Build Status](https://travis-ci.org/freedge/gomeme.svg?branch=master)](https://travis-ci.org/freedge/gomeme)
[![Build Status](https://dev.azure.com/freedge/freedge/_apis/build/status/freedge.gomeme?branchName=master)](https://dev.azure.com/freedge/freedge/_build/latest?definitionId=1&branchName=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/freedge/gomeme)](https://goreportcard.com/report/github.com/freedge/gomeme)
[![GoDoc](https://godoc.org/github.com/freedge/gomeme?status.svg)](https://godoc.org/github.com/freedge/gomeme) 


# gomeme



gomeme is a cli for Control-M, based on
[Control-M Automation API](https://docs.bmc.com/docs/automation-api/9181/services-784100995.html)
and loosely inspired by [govc](https://github.com/vmware/govmomi/tree/master/govc).

## Usage

```
export GOMEME_ENDPOINT=https://.../automation-api
export GOMEME_INSECURE=true
gomeme login -u toto
gomeme qr -n PRD*

```

Usage on Windows:
```
$env:GOMEME_INSECURE="true"
...
```

Use the ```--dump``` option to output in a Go-like format. ```--json``` option outputs in json format.

Helps for commands can be accessed with ```-h```

Traffic received from ctm server can be dumped with ```--debug```

## Commands

### login

Get a token for a user. Writes it into a .token file in the current directory.

The password must be provided either through the GOMEME_PASSWORD environment variable, or
through the terminal.

User defaults to $USER if not specified.

```
gomeme login --user toto
```

### qr

list qrs

```
gomeme qr --name PRD-*
```

### lj

list jobs (default limit is 1000). Use ```-v``` for more info. Outputs a csv with ```-csv```

```
gomeme lj --application TOTO-PRD --status Executing --limit 30
gomeme lj -a TOTO-PRD --host *pk1*
```

when listing a single jobid, one can use the deps option to go through the jobs in the neighbour of that job.

```
gomeme lj --deps --jobid FOOSRV:5rxwz -v
```


### qr.set

set a qr

```
gomeme qr.set --name DEV-FOO --ctm BARCT4T --max 5
```

### job.log

get the output or logs of a job id

```
gomeme job.log --jobid FOOCCT4P:5nq1c
gomeme job.log --jobid FOOCCT4P:5nq1c --output
gomeme job.log --jobid FOOCCT4P:5nq1c --output --run 3
```

### job.order

order a job.

```
gomeme job.order --ctm FOOCT4T --folder ABC-DEV-OND --jobs dABC1
```

### job.action

hold/delete/undelete/confirm/setToOk/rerun a job

```
gomeme job.action --action delete --jobid FOOCT4T:3z553
```

### job.tree

Tries to draw a tree of jobs with their dependencies.
This generates quite a lot of queries to control-m (1 per job to retrieve dependencies),
so it must be run carefully. Ensure there is no more than 100 jobs to analyse
when running tree.

The output looks like
```
A
  B
    C
  D
E
  F
```
which means that A and E do not have known predecessor
(in the list of selected nodes) and that B must complete before C starts, A must complete
before B and D starts.

We ensure the longest chain of jobs is shown, so if all jobs take the same amount of time to complete, in above
example C should start after B finishes, but there could be dependencies not appearing, such as C depending on
the completion of D, that are not reflected.

This command takes the same parameters as lj

```
gomeme job.tree --application TOTO-PRD --limit 10
gomeme job.tree --application TOTO-PRD --limit 10 --back    # instead of following dependencies, follow predecessors
```


### curl

Just outputs the curl command to run to target the API by hand

```
gomeme curl
```

### config.servers

List all servers

```
gomeme config.servers
```

### config.agents

List all agents for a specific server

```
gomeme config.servers --ctm FOO123
```


### config.agent

List parameters specific to an agent. Uses ```-all``` to show all parameters.

```
gomeme config.server --ctm FOO123 --host toto.net
```

### config.ping

Ping an agent, so that it ends up in config.agents output. Default timeout is 10s, use d to discover.

```
gomeme config.ping -H foo -c FOO
```

### deploy.get

Return definition of jobs in folder

```
gomeme deploy.get --ctm FOO --folder toto --xml
gomeme deploy.get --ctm FOO --folder toto
```

### deploy.put

🚧 Upload the definition of jobs

```
gomeme deploy.put --filename foo.json
```



### logout

Logout

```
gomeme logout
```


## License

http://www.apache.org/licenses/LICENSE-2.0
