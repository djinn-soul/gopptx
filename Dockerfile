# syntax=docker/dockerfile:1
FROM golang:1.25.7

RUN useradd -m -u 1000 appuser

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .
RUN chown -R appuser:appuser /workspace

USER appuser
CMD ["go", "test", "./..."]
