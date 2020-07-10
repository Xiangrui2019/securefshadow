FROM golang:latest as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go mod vendor
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o socks ./cmd/server

FROM ubuntu:latest
LABEL MAINTAINER="homeserving <xiangrui@aiursoft.com>"
COPY --from builder /app/socks .
EXPOSE 7748
CMD ./socks


