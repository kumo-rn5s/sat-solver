FROM golang:1.17.8-alpine3.15 as builder
WORKDIR /build
RUN apk add make
COPY . .
RUN make build

FROM alpine:3.15
COPY --from=builder /build/sat-solver /build/integration-test.sh /build/time-test.sh ./

ENTRYPOINT ["./sat-solver"]
