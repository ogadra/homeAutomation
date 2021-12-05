FROM golang:latest

RUN mkdir /go/src/work
WORKDIR /go/src/work

COPY . .

ENV GO111MODULE=on
RUN go get -u github.com/cosmtrek/air
EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
