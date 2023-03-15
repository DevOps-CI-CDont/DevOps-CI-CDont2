# Minitwit
A simple Twitter clone written in Next.js using the Go:Gin API framework.

## Getting Started
In order to run this project, you will need docker, to build and compose the images.

The backend is build...
```zsh
docker ... cicdont/backend
```

The frontend is build using the following command:
```zsh
docker ... cicdont/frontend
```

The next.js project is setup to use hotreloading. Run the frontend using the following command:
```zsh
docker-compose up -d --build
```