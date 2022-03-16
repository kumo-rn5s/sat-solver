FROM golang:1.17.8-alpine3.15 as builder
WORKDIR /app
COPY . .
RUN apk add make && make build
CMD ["./sat-solver"]
