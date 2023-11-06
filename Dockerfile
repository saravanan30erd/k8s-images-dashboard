# build stage
FROM golang:1.20 AS build-env

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV CGO_LDFLAGS="-s -w"
ENV GO111MODULE=on

COPY . /go/src/app
WORKDIR /go/src/app

RUN go build -o app

# final stage
FROM alpine:latest

RUN addgroup --system app \
    && adduser --system --ingroup app app

USER app

WORKDIR /app/
COPY --from=build-env /go/src/app /app/

RUN chown app /app/

CMD ./app
