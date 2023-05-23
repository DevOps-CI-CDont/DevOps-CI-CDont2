#!/bin/bash
echo "deleting manager1"
doctl compute droplet delete manager1 --force
# remove domain record pointing to manager1
doctl compute domain records delete cicdont.live 1 --force