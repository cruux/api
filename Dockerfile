FROM golang:latest

RUN go get github.com/cespare/reflex

ADD . /go/src/github.com/cruux/api
WORKDIR /go/src/github.com/cruux/api

RUN go install
CMD ["reflex", "-r", "\\.go$", "-s", "--", "sh", "-c", "go build && ./api"]
