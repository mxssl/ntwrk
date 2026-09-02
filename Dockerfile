# syntax=docker/dockerfile:1

FROM golang:1.27.1-alpine3.23@sha256:d9e2f2f07b10cc922da3e80e035c3058810b328d5aef82d2c63680967c5e2ec9 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
	go build -trimpath -ldflags="-s -w" -o /out/ntwrk .

FROM gcr.io/distroless/static-debian13:nonroot@sha256:f7f8f729987ad0fdf6b05eeeae94b26e6a0f613bdf46feea7fc40f7bd72953e6

COPY --from=builder --chown=65532:65532 --chmod=0555 /out/ntwrk /ntwrk

USER nonroot:nonroot
ENTRYPOINT ["/ntwrk"]
