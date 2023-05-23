#!/bin/bash
echo "deleting manager1"
doctl compute droplet delete manager1 --force
