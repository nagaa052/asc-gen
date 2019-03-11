FROM golang:alpine
MAINTAINER nagaa052
RUN apk --no-cache add ca-certificates

ADD build/asc-gen_unix /asc-gen

ENTRYPOINT ["/asc-gen"]
