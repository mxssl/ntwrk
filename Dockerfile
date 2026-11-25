# syntax=docker/dockerfile:1

FROM golang:1.27.4-alpine3.23@sha256:f70850ef6ce341e8b6427f677049145d63c0b7ce002d4d7e3d464828b1b96dcd AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
	go build -trimpath -ldflags="-s -w" -o /out/ntwrk .

FROM gcr.io/distroless/static-debian13:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7

COPY --from=builder --chown=65532:65532 --chmod=0555 /out/ntwrk /ntwrk

USER nonroot:nonroot
ENTRYPOINT ["/ntwrk"]
