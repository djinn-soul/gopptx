# syntax=docker/dockerfile:1
FROM golang:1.25.7

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .

CMD ["go", "test", "./..."]
