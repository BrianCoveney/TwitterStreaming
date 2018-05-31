FROM golang:1.9.3

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build repository-service.go

EXPOSE 9090
ENTRYPOINT ["/go/repository-service"]