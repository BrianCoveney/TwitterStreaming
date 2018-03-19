FROM golang:1.9

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build server.go

ENTRYPOINT ["/go/server"]
