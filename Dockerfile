#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/main -v ./cmd/api/main.go

#final stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/main /go/bin/.en.vault /app/
ENTRYPOINT /app/main
LABEL Name=goweather Version=0.0.1
EXPOSE 8003
