FROM golang:1.16-alpine as builder

RUN mkdir /build
WORKDIR /build

COPY . .

RUN export GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tcp-echo .

FROM alpine:latest

ADD ./tls /build/tls
COPY --from=builder /build/tcp-echo ./tcp-echo

WORKDIR /

ENTRYPOINT ["/tcp-echo"]