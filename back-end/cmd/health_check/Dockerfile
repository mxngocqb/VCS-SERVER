FROM golang:1.22

WORKDIR /go/app

COPY . .

RUN apt-get update && apt-get install -y iputils-ping
RUN go mod tidy
RUN go build ./cmd/health_check

ENTRYPOINT [ "./health_check"]