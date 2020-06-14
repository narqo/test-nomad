# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV["LC_ALL"] = "en_US.UTF-8"

$num_a_nodes = 3
$num_b_nodes = 2

Vagrant.configure("2") do |config|
  #config.vm.box = "bento/ubuntu-18.04"
  config.vm.box = "debian/buster64"
  config.vm.box_check_update = false

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "512"
  end

  # forward ssh agent to easily ssh into the different machines
  config.ssh.forward_agent = true
  # always use Vagrants insecure key
  config.ssh.insert_key = false

  config.vm.network "forwarded_port", guest: 8080, host: 8080, auto_correct: true, host_ip: "127.0.0.1"

  config.vm.provision "ansible" do |ansible|
    ansible.playbook = "site.yml"
    ansible.inventory_path = "inventory/10-hosts"
  end

  (1..$num_a_nodes).each do |n|
    config.vm.define "dc-a-node-#{n}" do |node|
      node.vm.hostname = "dc-a-node-#{n}"
      node.vm.network "private_network", ip: "172.16.1.#{n+100}"
      if n == 1
        # Expose nomad ports for API and UI
        node.vm.network "forwarded_port", guest: 4646, host: 4646, auto_correct: true, host_ip: "127.0.0.1"
      end
    end
  end

  (1..$num_b_nodes).each do |n|
    config.vm.define "dc-b-node-#{n}" do |node|
      node.vm.hostname = "dc-b-node-#{n}"
      node.vm.network "private_network", ip: "172.16.2.#{n+100}"
    end
  end
end
