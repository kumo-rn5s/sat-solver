FROM golang:1.17.8-alpine3.15 as builder
WORKDIR /build
COPY . /build/
RUN GOOS=linux GOARCH=arm64 go build

FROM alpine:3.15
COPY --from=builder /build/sat-solver .

ENTRYPOINT ["/bin/sh","-c","./sat-solver"]