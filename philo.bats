#!/usr/bin/env bats

@test "login" {
  gomeme login -u workbench
}

@test "define our philosophers and chopsticks" {
  gomeme deploy.put --subject test -f commands/deploy/fixtures/philosophers.json -c workbench >&3
  for name in a b c d e ; do
    gomeme test.qr.new --subject test -c workbench -m 1 -n $name
  done
  sleep 1
  gomeme qr >&3
}

@test "get them to eat" {
  for name in 1 2 3 4 5 ; do
    gomeme job.order --subject test -c workbench -f philo -n $name -D >&3
  done
}

@test "check who ate" {
  for i in `seq 20` ; do
    gomeme lj -c workbench -a appli -v >&3
    sleep 10
  done
  for name in 1 2 3 4 5 ; do 
    ID=$(gomeme lj -n $name --json | jq '.Statuses[0].JobId')
    gomeme job.log -j ${ID} >&3
    gomeme job.log -j ${ID} -o >&3
    gomeme job.get -j ${ID} >&3
  done
}