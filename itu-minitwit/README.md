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

## Testing locally

to test the backend locally, you need to run the backend via `go run main.go` from the backend/Main directory. When it is running you can run the test via `go test` from the backend/test directory.

## Running with docker

`docker compose -f docker-compose-dev.yml up` from itu-minitwit directory