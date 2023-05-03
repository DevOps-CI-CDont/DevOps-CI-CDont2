#!/bin/bash
echo "deleting manager1, worker1, worker2"
doctl compute droplet delete manager1 worker1 worker2 --force
