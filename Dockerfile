FROM golang:1.17.8-alpine3.15
WORKDIR /app
COPY . .
RUN apk add make && make build
CMD ["./sat-solver"]
