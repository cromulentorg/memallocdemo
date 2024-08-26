FROM golang:1.21-alpine as builder

ENV BUILD_DIR /build

COPY src $BUILD_DIR/src
WORKDIR $BUILD_DIR/src


RUN go get
RUN go test ./...
RUN go build -o /dist/memallocdemo

FROM alpine

RUN apk update \
  && apk upgrade \
  && apk add ca-certificates

COPY --from=builder /dist/memallocdemo /bin/

# create least privileged user and use it
RUN adduser --disabled-password --gecos '' user-app
USER user-app

EXPOSE 8080
EXPOSE 9090
CMD ["memallocdemo"]
