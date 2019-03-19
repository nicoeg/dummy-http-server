FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o dummyhttpserver .

FROM scratch

COPY --from=builder /build/dummyhttpserver /app/
COPY config.json /app/config.json

EXPOSE 8080

WORKDIR /app

ENTRYPOINT ["/app/dummyhttpserver"]
