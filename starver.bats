#!/usr/bin/env bats

@test "login" {
  gomeme login -u workbench
}

@test "define our philosophers and chopsticks" {
  gomeme deploy.put --subject test -f commands/deploy/fixtures/philosophersstarve.json -c workbench >&3
  for name in a b c d e f g h i j k l m; do
    gomeme test.qr.new --subject test -c workbench -m 1 -n $name >&3 || echo .
  done
  sleep 1
  gomeme qr >&3
}

@test "get them to eat" {
  for name in `seq 13` ; do
    gomeme job.order --subject test -c workbench -f philo2 -n $name -D >&3
  done
   gomeme job.order --subject test -c workbench -f philo2 -n starver -D >&3
}

@test "check who ate" {
  gomeme lj -c workbench -a appli -v >&3
  for name in `seq 13` ; do 
    ID=$(gomeme lj -n $name --json | jq '.Statuses[0].JobId')
    gomeme job.get -j ${ID} >&3
  done
}
