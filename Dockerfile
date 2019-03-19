FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o dummy-http-server .

FROM scratch

COPY --from=builder /build/dummy-http-server /app/
COPY config.json /app/config.json

EXPOSE 8080

WORKDIR /app

ENTRYPOINT ["/app/dummy-http-server"]
