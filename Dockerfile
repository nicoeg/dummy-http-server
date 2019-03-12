FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o dummy-http-server .

FROM scratch

COPY --from=builder /build/dummy-http-server /app/

EXPOSE 8080

ENTRYPOINT ["/app/dummy-http-server"]
