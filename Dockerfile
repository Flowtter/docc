FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /docc

FROM alpine
COPY --from=builder /docc /docc
COPY --from=builder /app/assets /assets

ENTRYPOINT /docc