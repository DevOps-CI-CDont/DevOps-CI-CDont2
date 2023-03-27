FROM golang:1.20.0
WORKDIR /app

COPY . .

CMD ["go", "run", "main.go", "test"]