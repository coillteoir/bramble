FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /lugh

LABEL stage=builder

#Ubuntu base image to test pod
FROM ubuntu:latest

WORKDIR /

COPY --from=builder /lugh /lugh

CMD ["/lugh"]
