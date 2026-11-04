# syntax=docker/dockerfile:1

FROM golang:1.27.3-alpine3.23@sha256:8a5046497cda39860b832ffdbe8dd66018a0f05bb01352f533a36e7e2ced66a0 AS builder

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
