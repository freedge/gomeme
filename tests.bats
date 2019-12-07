#!/usr/bin/env bats

@test "login" {
  run gomeme login -u workbench
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Logged in. Server version " ]]
}

@test "curl" {
  gomeme curl  
}

@test "bootstrap" {
  gomeme test.config.emparamset --subject test --description test --name UserAuditAnnotationOn --value 1
  gomeme test.qr.new --subject test --debug -n INIT -m 0 -c workbench
}

@test "qr" {
  gomeme --debug qr >&3
}

@test "config server should retrieve our workbench server" {
  gomeme config.servers >&3
  gomeme config.agents -c workbench >&3
}

@test "ping our agent and show its parameters" {
  until gomeme config.ping  -c workbench --host workbench -t 10 ; do sleep 1 ; done
  gomeme config.agent -c workbench --host workbench   
}

@test "deploy.put should put some jobs" {
  run gomeme deploy.put --subject test -f commands/deploy/fixtures/folder.json -c workbench
  [ "$status" -eq 0 ]
  [[ "$output" =~ "1 jobs successfully deployed" ]]
}

@test "deploy.get returns the defined jobs" {
  run gomeme deploy.get -c workbench -f FOO-BARLOCAL-PRK
  [ "$status" -eq 0 ]
  DESC=$(echo "$output" | jq '."FOO-BARLOCAL-PRK".dFOOJOBPRGPK1.Description')
  [ "$DESC" == '"1234"' ]
}

@test "deploy.get --xml returns the defined jobs" {
  gomeme deploy.get -c workbench -f FOO-BARLOCAL-PRK --xml | xmllint --encode ASCII - | grep CREATED_BY >&3  
}

@test "job.order to order them and keep them held" {
  run gomeme job.order --subject test -c workbench -f FOO-BARLOCAL-PRK -n dFOOJOBPRGPK1
  [ "$status" -eq 0 ]
  [[ "${lines[0]}" =~ "RunId:" ]]
  [[ "${lines[1]}" =~ "JobId:" ]]
}

@test "list the jobs" {
  # the ordered job is not available immediately
  until gomeme lj --json | jq .Statuses[0] ; do echo . ; sleep 1 ; done
  gomeme lj -H '*BA*'
  gomeme lj -n dFOOJOBPRGPK1
  gomeme lj --debug -c workbench >&3
  run gomeme lj -v
  [[ "$output" =~ "1/1" ]]
  [ "$status" -eq 0 ]
}

@test "list the jobs as a tree" {
  run gomeme job.tree -v -H '*BA*'
  [ "$status" -eq 0 ]
  [[ "$output" =~ dFOOJOBPRGPK1 ]]
}

@test "list the jobs as csv" {
  run gomeme lj --csv
  [ "$status" -eq 0 ]
  [[ "${lines[1]}" == "FOO-BARLOCAL-PRK,dFOOJOBPRGPK1,Wait Condition,-1.000000,," ]]
}

@test "get the jobs and its dependencies" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.tree --deps -v -j ${ID}
  [[ "$output" =~ dFOOJOBPRGPK1 ]]
  [ "$status" -eq 0 ]
}  

@test "change the job parameters" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.modify -j $ID -n dFOOJOBPRGPK1 --subject test --debug -- A B C
  echo "$output" >&3
  [ "$status" -eq 0 ]
  [[ "$output" =~ "job was successfully modifyed" ]]
}  

@test "free the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a free --subject test
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully freed" ]]
}  

@test "get the job waiting info" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  gomeme job.get --debug -j $ID >&3
  run gomeme job.get -j $ID  
  [[ "$output" =~ "There is no machine available for job execution" ]]
} 

@test "hold the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a hold --subject test
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully held" ]]
}

@test "get the job logs" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  gomeme job.log -j $ID >&3
}  

@test "delete the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a delete --subject test
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully deleted" ]]
}

@test "setting QR" { 
  gomeme qr.set -n INIT -c workbench -m 42 --subject test
  QR=$(gomeme qr -n INIT --json | jq '.[0].Max')
  until [ "$QR" -eq 42 ] ; do
    echo ${QR} >&3
    sleep 1
    QR=$(gomeme qr -n INIT --json | jq '.[0].Max')
  done
}

@test "create a secret and check for its existence" { 
  echo 42 | gomeme secret.add --subject test -n mysecret -f -
  echo 好 | gomeme secret.update --subject test -n mysecret -f -
  gomeme secret.get | grep mysecret
}

@test "logout" {
  run gomeme logout
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Successfully logged out from session" ]]
}
