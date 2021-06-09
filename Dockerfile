#
#  Build Container
#
FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /src

RUN apk add build-base git

COPY go.sum go.mod /src/

WORKDIR /src

COPY . /src

RUN make install && \
    mkdir -p /app && \
    cp -r ${GOPATH}/bin/urlshortener /app/

#
# Runtime Container
#
FROM alpine:3.9

ENV PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    ca-certificates bash && \
    mkdir -p /var/log && \
    chgrp -R 0 /var/log && \
    chmod -R g=u /var/log

WORKDIR /app

COPY --from=builder /app /app/

COPY --from=builder /src/docker/urlshortener/run.sh /app/run.sh

CMD /app/run.sh
