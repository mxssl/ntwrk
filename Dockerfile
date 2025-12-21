FROM golang:1.25-alpine3.22 as builder

WORKDIR /go/src/github.com/mxssl/ntwrk
COPY . .

# Install external dependcies
RUN apk add --no-cache \
	ca-certificates \
	curl \
	git

# Compile binary
RUN CGO_ENABLED=0 \
	go build -v -o ntwrk

# Copy compiled binary to clear Alpine Linux image
FROM alpine:3.23.0
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mxssl/ntwrk/ntwrk /ntwrk
RUN chmod +x ntwrk
CMD ["./ntwrk"]
