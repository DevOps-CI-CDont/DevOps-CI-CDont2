#!/bin/bash
echo "deleting manager1"
doctl compute droplet delete manager1 --force
# remove domain record pointing to manager1
# Get all DNS A records for the domain
records=$(doctl compute domain records list "cicdont.live" --format "ID,Type" --no-header | grep "A" | awk '{print $1}')
for record_id in $records; do
  doctl compute domain records delete "cicdont.live" "$record_id"
done