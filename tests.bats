#!/usr/bin/env bats

@test "login" {
  run gomeme login -u workbench
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Logged in. Server version " ]]
}

@test "bootstrap" {
  run gomeme config.emparamset --debug --subject test --name UserAuditAnnotationOn --value 1
  # ignoring the output of this for the time being
  echo "[$status] $output" >&3
}

@test "config server should retrieve our workbench server" {
  gomeme config.servers >&3
  gomeme config.agents -c workbench >&3
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

@test "job.order to order them and keep them held" {
  run gomeme job.order -c workbench -f FOO-BARLOCAL-PRK -n dFOOJOBPRGPK1
  [ "$status" -eq 0 ]
  [[ "${lines[0]}" =~ "RunId:" ]]
  [[ "${lines[1]}" =~ "JobId:" ]]
}

@test "list the jobs" {
  # the ordered job is not available immediately
  until gomeme lj --json | jq .Statuses[0] ; do echo . ; sleep 1 ; done
  gomeme lj -H '*BA*' >&3
  gomeme lj -n dFOOJOBPRGPK1 >&3
  gomeme lj -c workbench >&3
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
  run gomeme job.modify -j $ID -n dFOOJOBPRGPK1 --subject test -- A B C
  [ "$status" -eq 0 ]
  [[ "$output" =~ "job was successfully modifyed" ]]
}  

@test "free the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a free
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully freed" ]]
}  

@test "get the job waiting info" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.get -j $ID
  [ "$status" -eq 0 ]
  [[ "$output" =~ "There is no machine available for job execution" ]]
} 

@test "hold the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a hold
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully held" ]]
}

@test "get the job logs" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  gomeme job.log -j $ID >&3
}  

@test "delete the job" {
  ID=$(gomeme lj --json | jq '.Statuses[0].JobId')
  run gomeme job.action -j $ID -a delete
  [ "$status" -eq 0 ]
  [[ "$output" =~ "was successfully deleted" ]]
}

@test "logout" {
  run gomeme logout
  [ "$status" -eq 0 ]
  [[ "$output" =~ "Successfully logged out from session" ]]
}