FROM golang:1.22

WORKDIR /go/app

COPY . .

RUN apt-get update
RUN go mod tidy
RUN go build ./cmd/api

ENTRYPOINT [ "./api"]