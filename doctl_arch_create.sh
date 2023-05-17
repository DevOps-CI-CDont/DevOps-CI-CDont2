#!/bin/bash
if doctl compute droplet list --format "Name" | grep "manager1"; then
    echo "Error: Droplet with name manager1 already exists."
    exit 1
fi


# define fingerprints
ssh_key_fingerprints="ea:e6:b9:88:5a:8f:d6:4b:c7:03:06:c2:fd:aa:ca:a8,3d:83:40:84:6a:ee:3b:e4:da:e6:b8:c9:1e:4f:24:e5,58:15:a0:48:24:47:dc:79:af:80:1b:f9:e3:0e:ff:67,cd:82:77:0a:e5:c3:b4:b2:06:91:30:b6:b3:60:f0:76,ee:36:4b:17:af:40:f4:0c:00:8e:bf:a9:46:ab:df:e4"
echo "ssh_key_fingerprints: $ssh_key_fingerprints"

env_file_path="itu-minitwit/backend/.env"

# check if file exists
if [ ! -f "$env_file_path" ]; then
    echo "Error: File $env_file_path does not exist."
    exit 1
fi

# create droplets
echo "creating droplets"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-2vcpu-4gb-amd --region fra1 --enable-monitoring manager1 --ssh-keys "$ssh_key_fingerprints"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-1vcpu-2gb --region fra1 --enable-monitoring worker1 --ssh-keys "$ssh_key_fingerprints"
doctl compute droplet create --image ubuntu-22-10-x64 --size s-1vcpu-2gb --region fra1 --enable-monitoring worker2 --ssh-keys "$ssh_key_fingerprints"

# wait for droplets to be created
echo "waiting 60 seconds for droplets to be created"
sleep 10
echo "50"
sleep 10
echo "40"
sleep 10
echo "30"
sleep 10
echo "20"
sleep 10
echo "10"
sleep 10
echo "Done waiting"


# get droplet IP addresses
manager1_ip=$(doctl compute droplet get manager1 --format PublicIPv4 --no-header)
wanky=$(doctl compute droplet get worker1 --format PublicIPv4 --no-header)
worker2_ip=$(doctl compute droplet get worker2 --format PublicIPv4 --no-header)

# print IP addresses
echo "Manager1 IP address: $manager1_ip"
echo "Worker1 IP address: $wanky"
echo "Worker2 IP address: $worker2_ip"

# add droplets to known hosts
sleep 1
ssh-keyscan -H $manager1_ip >> ~/.ssh/known_hosts
sleep 1
ssh-keyscan -H $wanky >> ~/.ssh/known_hosts
sleep 1
ssh-keyscan -H $worker2_ip >> ~/.ssh/known_hosts

# Check that all dropslets are added to known hosts otherwise run again
# add droplets to known hosts
for ip in $manager1_ip $wanky $worker2_ip; do
    while ! ssh-keygen -F $ip | grep -q $ip; do
        echo "Adding $ip to known hosts failed...retrying"
        sleep 5
        ssh-keyscan -H $ip >> ~/.ssh/known_hosts
    done
done

# install docker engine on droplets
doctl compute ssh manager1 --ssh-command "curl -fsSL https://get.docker.com -o get-docker.sh"

# get tokens from .env file
echo "trying to get tokens from .env file"
source "$env_file_path"
echo "manager_token: $manager_token"
echo "worker1_token: $worker1_token"
echo "worker2_token: $worker2_token"



echo "installing doctl on droplet"
doctl compute ssh manager1 --ssh-command "sudo snap install doctl"
echo "manager1 doctl auth init with token"
doctl compute ssh manager1 --ssh-command "sudo mkdir /root/.config"
doctl compute ssh manager1 --ssh-command "sudo snap connect doctl:dot-docker"
doctl compute ssh manager1 --ssh-command "sudo doctl auth init -t $manager_token"
echo "manager1 doctl registry login"
doctl compute ssh manager1 --ssh-command "sudo doctl registry login"

# get docker compose file on manager 1
echo "trying to get docker compose file on manager1"
doctl compute ssh manager1 --ssh-command "curl https://raw.githubusercontent.com/DevOps-CI-CDont/DevOps-CI-CDont/IaC/docker-compose-manager.yml --output ./docker-compose.yml "
doctl compute ssh manager1 --ssh-command "curl https://raw.githubusercontent.com/DevOps-CI-CDont/DevOps-CI-CDont/IaC/itu-minitwit/nginx.conf --output /nginx.conf"
doctl compute ssh manager1 --ssh-command "curl https://raw.githubusercontent.com/DevOps-CI-CDont/DevOps-CI-CDont/IaC/itu-minitwit/.htpasswd --output /.htpasswd"

echo "docker compose up on manager1"
doctl compute ssh manager1 --ssh-command "sh get-docker.sh"
scp $env_file_path root@$manager1_ip:./.env
doctl compute ssh manager1 --ssh-command "docker compose up"
