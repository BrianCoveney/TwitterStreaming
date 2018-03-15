FROM golang:1.9

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN go get ./...

RUN go build twitter-service.go

ENTRYPOINT ["/go/twitter-service"]
