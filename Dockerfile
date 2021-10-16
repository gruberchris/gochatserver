FROM golang:1.17-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o gochatserver .

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app/ /app/
EXPOSE 5000
ENTRYPOINT ["./gochatserver"]