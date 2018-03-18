FROM golang:1.9.3

WORKDIR /go
COPY . .

RUN go get -d -v ./...
RUN go build front-end.go

EXPOSE 3000
ENTRYPOINT ["/go/frontendservice"]
