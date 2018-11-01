FROM golang:alpine as builder

WORKDIR /go/src/github.com/mxssl/ntwrk.cf
COPY . .

# Compile
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o app

# Copy compiled binary to clear Alpine Linux image
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mxssl/ntwrk.cf .
RUN chmod +x app
CMD ["./app"]
