FROM golang:1-alpine AS builder


WORKDIR /app

COPY *.go go.mod go.sum /app/

RUN go mod download
RUN GOEXPERIMENT=greenteagc go build -o vcspammer .


FROM alpine:latest


WORKDIR /app

RUN apk --no-cache add ffmpeg

COPY --from=builder /app/vcspammer .
RUN chmod +x vcspammer

ENTRYPOINT ./vcspammer
