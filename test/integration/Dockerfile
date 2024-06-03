FROM golang:1.22

ARG WORKDIR

WORKDIR $WORKDIR

COPY go.mod go.sum ./

RUN go mod download
