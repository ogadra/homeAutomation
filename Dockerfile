FROM golang:latest

RUN mkdir /go/src/work
WORKDIR /go/src/work

ADD . /go/src/work

ENV GO111MODULE=on
RUN go mod download
EXPOSE 8080
CMD ["go", "run", "/go/src/main.go"]