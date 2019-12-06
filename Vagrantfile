# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/xenial64" 
  # access a port on your host machine (via localhost) and have all data forwarded to a port on the guest machine.
  config.vm.network "forwarded_port", guest: 8443, host: 8443
  config.vm.hostname = "workbench"
  # Create a private network, which allows host-only access to the machine
  config.vm.provider "virtualbox" do |vb|
    vb.name = 'ctmworkbench'
    vb.memory = 4096
    vb.cpus = 1
  end
  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/freedge/gomeme"

  # get our workbench running
  config.vm.provision "docker" do |d|
    d.post_install_provision "shell", 
      inline: "docker images | grep controlm-workbench | grep 9.19.200 || ( curl -L https://controlm-appdev.s3-us-west-2.amazonaws.com/workbench/9.0.19.200/deploy/ova/controlm-workbench-9.19.200.xz | docker load )"
    d.run "workbench",
      image: "controlm-workbench:9.19.200",
      args: "--tty --hostname workbench -p 8443:8443 -p 7005:7005" 
  end

  # retrieve the certificate to use to connect to it
  config.vm.provision "shell",
    inline: "mkdir -p /home/vagrant/go/src/github.com/freedge/gomeme/.certs ; until echo | openssl s_client -prexit -connect workbench:8443 | openssl x509 > /home/vagrant/go/src/github.com/freedge/gomeme/.certs/out.pem ; do sleep 1 ; done",
    privileged: false
  # install the needed packages
  config.vm.provision "shell",
    inline: "snap install --classic go"
  config.vm.provision "shell",
    inline: "snap install jq"
  config.vm.provision "shell",
    inline: "snap install libxml2-utils"
  config.vm.provision "shell",
    inline: "cd ; rm -rf bats-core-1.1.0 ; curl -L https://github.com/bats-core/bats-core/archive/v1.1.0.tar.gz | tar zxf - ; bats-core-1.1.0/install.sh /usr"
  config.vm.provision "shell",
    inline: "echo 'cd /home/vagrant/go/src/github.com/freedge/gomeme ; export GOMEME_CERT_DIR=`pwd`/.certs ; export PATH=`pwd`:${PATH} ; export GOMEME_PASSWORD=workbench ; export GOMEME_ENDPOINT=https://workbench:8443/automation-api' >> /home/vagrant/.profile ",
    privileged: false
  config.vm.provision "shell",
    inline: "cd /home/vagrant/go/src/github.com/freedge/gomeme && go get ./... && go build",
    privileged: false


end
