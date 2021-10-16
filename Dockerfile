FROM golang:1.17-alpine

RUN mkdir /app
WORKDIR /app

RUN export GO111MODULE=on
ADD . .

RUN go build

EXPOSE 5000

ENTRYPOINT ["/app/gochatserver"]