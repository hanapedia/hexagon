FROM golang:1.23

ARG WORKDIR

WORKDIR $WORKDIR

COPY go.mod go.sum ./

RUN go mod download
