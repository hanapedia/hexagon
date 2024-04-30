FROM golang:1.21

ARG WORKDIR

WORKDIR $WORKDIR

COPY go.mod go.sum ./

RUN go mod download
