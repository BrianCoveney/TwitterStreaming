FROM golang:1.9

MAINTAINER brian.coveney

WORKDIR /go
COPY . .

#RUN bash get.sh

RUN go get ./...

RUN go build frontend.go

EXPOSE 3000
ENTRYPOINT ["/go/frontend"]
