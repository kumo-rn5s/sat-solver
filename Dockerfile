FROM golang:alpine as builder
RUN mkdir /build
ADD *.go /build/
ADD go.* /build/
WORKDIR /build
RUN go test && CGO_ENABLED=0 GOOS=linux go build

FROM alpine/make:latest
COPY --from=builder /build/sat-solver .
COPY ./test ./test
COPY Makefile .

ENTRYPOINT [ "sh","-c","make integration-test" ]