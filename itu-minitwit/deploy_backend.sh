doctl registry login
docker tag backend registry.digitalocean.com/cicdont-images/backend
docker push registry.digitalocean.com/cicdont-images/backend