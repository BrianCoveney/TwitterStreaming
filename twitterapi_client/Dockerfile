FROM golang:1.9.3

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build twitterapi_client.go

ENTRYPOINT ["/go/twitterapi_client"]
