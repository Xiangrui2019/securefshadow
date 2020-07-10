FROM golang:latest as builder
WORKDIR /app
COPY go.mod go.sum ./
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=gosum.io+ce6e7565+AY5qEHUk/qmHc5btzW45JVoENfazw8LielDsaI+lEbq6
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o socks ./cmd/server

FROM ubuntu:latest
LABEL MAINTAINER="homeserving <xiangrui@aiursoft.com>"
COPY --from=builder /app/socks .
EXPOSE 7748
CMD ./socks


