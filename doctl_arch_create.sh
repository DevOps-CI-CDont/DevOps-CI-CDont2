#!/bin/bash
# fingerprints
ssh_key_fingerprints="ea:e6:b9:88:5a:8f:d6:4b:c7:03:06:c2:fd:aa:ca:a8,3d:83:40:84:6a:ee:3b:e4:da:e6:b8:c9:1e:4f:24:e5,58:15:a0:48:24:47:dc:79:af:80:1b:f9:e3:0e:ff:67,cd:82:77:0a:e5:c3:b4:b2:06:91:30:b6:b3:60:f0:76"
echo "ssh_key_fingerprints: $ssh_key_fingerprints"

# create droplets
echo "creating droplets"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-2vcpu-4gb-amd --region fra1 --enable-monitoring manager1 --ssh-keys "$ssh_key_fingerprints"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-1vcpu-2gb --region fra1 --enable-monitoring worker1 --ssh-keys "$ssh_key_fingerprints"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-1vcpu-2gb --region fra1 --enable-monitoring worker2 --ssh-keys "$ssh_key_fingerprints"

# wait for droplets to be created
echo "waiting 60 seconds for droplets to be created"
sleep 60

# install docker engine on manager1
echo "trying to install docker engine on manager1"
doctl compute ssh manager1 --ssh-command "sudo snap install docker"
echo "docker engine installed on manager1"
echo "trying to install docker engine on worker1"
doctl compute ssh worker1 --ssh-command "sudo snap install docker"
echo "docker engine installed on worker1"
echo "trying to install docker engine on worker2"
doctl compute ssh worker2 --ssh-command "sudo snap install docker"
echo "docker engine installed on worker2"

# get docker compose file on manager 1
echo "trying to get docker compose file on manager1"
doctl compute ssh manager1 --ssh-command "curl https://raw.githubusercontent.com/DevOps-CI-CDont/DevOps-CI-CDont/main/itu-minitwit/docker-compose.yml --output ./docker-compose.yml "

echo "docker compose up on manager1"
doctl compute ssh manager1 --ssh-command "cd DevOps-CI-CDont/itu-minitwit\
&& docker compose up"

