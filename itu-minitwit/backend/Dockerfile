FROM golang:1.20 AS builder

WORKDIR /user/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

WORKDIR /user/src/app
RUN go build -v -o /usr/local/bin/app ./Main

EXPOSE 8080
EXPOSE 8081

ENV PORT 8080
ENV API_BASE_URL http://138.68.93.147:8080

CMD ["app"]