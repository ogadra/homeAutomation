FROM golang:latest

RUN mkdir /go/src/work
WORKDIR /go/src/work

COPY . .

ENV GO111MODULE=on
RUN go mod download
EXPOSE 8080

CMD ["go", "run", "main.go"]
