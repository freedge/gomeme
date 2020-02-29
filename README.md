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
go install github.com/freedge/gomeme

export GOMEME_ENDPOINT=https://.../automation-api
export GOMEME_CERT_DIR=~/certs
gomeme login -u toto
gomeme qr -n PRD*
...
```


Use the ```--dump``` option to output in a Go-like format. ```--json``` option outputs in json format.

Helps for commands can be accessed with ```-h```

Traffic received from ctm server can be dumped with ```--debug```

Usage on Windows:

It might be installed using [the latest release](https://github.com/freedge/gomeme/releases/latest/download/gomeme.appinstaller)

```
$env:GOMEME_ENDPOINT="https://workbench:8443/automation-api"
$env:GOMEME_CERTDIR=...
gomeme login -u workbench

$Jobs = & gomeme lj -s Executing -f *appli* --json | ConvertFrom-Json
$Jobs.statuses |  Out-GridView
```


### SSL config

If GOMEME_ENDPOINT is using https, gomeme must be able to verify the identity of the server.
When the certificate of the server is signed by an unknown authority, that certificate should be placed under a folder 
referenced by SSL_CERT_DIR or GOMEME_CERT_DIR environment variable.

Typically,

1. create a certs folder
2. retrieve the server certificate using ```echo | openssl s_client -prexit -connect myserver:443 | openssl x509 > certs/out.pem```
3. ensure the server name in the certificate can be resolved (possibly, add it in etc/hosts)
4. from the certs folder, run ```c_rehash .``` (so that curl will be able to use --capath, though it does not seem mandatory)

SSL_CERT_DIR usually can be a column separated list of folders, but go (https://github.com/golang/go/issues/35325) consider it
as a single folder. GOMEME_CERT_DIR environment variable can be used as an alternative to SSL_CERT_DIR (to avoid breaking
other tools).

openssl and go tls package do not validate certificates identically (it seems you need to trust the whole chain, including
the root CA, for openssl to work fine). As such, you may want to run 
```echo | openssl s_client -prexit -connect myserver:443 -showcerts``` 
and put the root CA certificate in your certs folder.

### Annotations

Annotations can be provided with the ```--subject``` and ```--description``` parameters. 
```deploy.put```, ```job.modify```, ```job.order```
require an annotation to be set even if audit is not activated on server side.


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

list jobs (default limit is 1000). Use ```-v``` for more info. Outputs a csv with ```--csv```

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
gomeme qr.set --name DEV-FOO --ctm BARCT4T --max 5 --subject reason
```

### job.log

get the output or logs of a job id

```
gomeme job.log --jobid FOOCCT4P:5nq1c
gomeme job.log --jobid FOOCCT4P:5nq1c --output
gomeme job.log --jobid FOOCCT4P:5nq1c --output --run 3
```

### job.order

order a job or a whole folder

```
gomeme job.order --ctm FOOCT4T --folder ABC-DEV-OND --subject Test1234
gomeme job.order --ctm FOOCT4T --folder ABC-DEV-OND --jobs dABC1 --subject Test1234
```

By default the job is held unless the ```-D``` option is provided.

Gomeme tries to get the job id of the created job and retry a few times (waiting 1s between each try). Default is 2 tries, use ```--retries 0``` to not wait.

An annotation is required.

### job.action

hold/delete/undelete/confirm/setToOk a job

```
gomeme job.action --action delete --jobid FOOCT4T:3z553
```

### job.rerun

rerun a job

```
gomeme job.rerun --jobid FOOCT4T:3z553
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

### ps

Same as curl but with a call to Invoke-RestMethod PowerShell cmdlet.

Powershell has no equivalent for capath, so we will skip server verification entirely in that case. Tested against PowerShell 7 preview.

```
$ps=& gomeme ps
Invoke-Expression $ps/config/servers
```

### config.servers

List all servers

```
gomeme config.servers
```

### config.server

List all agents and hostgroups for a specific server

```
gomeme config.server --ctm FOO123
```


### config.agent

List parameters specific to an agent. Uses ```--all``` to show all parameters.

```
gomeme config.agent --ctm FOO123 --host toto.net
```

### config.ping

Ping an agent, so that it ends up in config.agents output. Default timeout is 10s, use d to discover.

```
gomeme config.ping -H foo -c FOO
```

### config.hostgroup

Returns the list of agents parts of a hostgroup

```
gomeme config.hostgroup -c foo -g group
```

### deploy.get

Return definition of jobs in folder

```
gomeme deploy.get --ctm FOO --folder toto --xml
gomeme deploy.get --ctm FOO --folder toto
```

### deploy.put

Upload the definition of jobs

```
gomeme deploy.put --filename foo.json --subject "Record1234" --ctm FOOCTM
```

### job.get

Retrieve status, waiting info, parameters of a single job

```
gomeme job.get -j FOO:123
```

### job.modify

Change the parameters of a held job

```
gomeme job.modify  --subject "Record 1234" -n jobname -j FOOBAR:09otj -- param1 param2
```

### secret

secret things

### logout

Logout

```
gomeme logout
```

## Local dev

```
vagrant plugin install vagrant-disksize
vagrant up
vagrant ssh
bats tests.bats
bats philo.bats
```

```vagrant up``` will start a localdev environment, including a running controlm workbench.

After adding
```
127.0.0.1 workbench
```

the workbench can be access from https://workbench:8443/automation-api

gomeme compiled under windows can be launched from VSCode terminal with
```
$env:GOMEME_ENDPOINT="https://workbench:8443/automation-api"
$env:GOMEME_PASSWORD="workbench"
$env:GOMEME_CERT_DIR=".certs"
.\gomeme.exe login -u workbench
```

however the bats tests must be started under Linux.

### Azure pipelines

Set-up for go:

https://docs.microsoft.com/en-us/azure/devops/pipelines/ecosystems/go?view=azure-devops&tabs=go-current

Create the manifest file:
https://docs.microsoft.com/en-us/windows/msix/desktop/desktop-to-uwp-manual-conversion

not sure if needed : Create the package manifest using MakeAppx.exe 

https://developer.microsoft.com/fr-fr/windows/downloads/windows-10-sdk/


Create a certificate 

https://docs.microsoft.com/en-us/windows/msix/package/create-certificate-package-signing

```
 New-SelfSignedCertificate -Type Custom -Subject "CN=Frigocorp, O=Frigocorp, L=Gotham, S=Hyrule, C=France" -KeyUsage DigitalSignature -FriendlyName "Frigo certificate" -CertStoreLocation "Cert:\CurrentUser\My" -TextExtension @("2.5.29.37={text}1.3.6.1.5.5.7.3.3", "2.5.29.19={text}")
```


Set-up the pipeline 

https://docs.microsoft.com/en-us/windows/msix/desktop/azure-dev-ops

## License

http://www.apache.org/licenses/LICENSE-2.0
