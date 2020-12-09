FROM golang:alpine AS builder

ARG ARCH=amd64

COPY go.mod /
COPY go.sum /

RUN go mod download

COPY main.go /
COPY viewproxy /viewproxy

# Ensure we don't need glibc on the target with these flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -ldflags="-w -s" /main.go

FROM scratch

COPY --from=builder /go/main /prometheus-view-proxy

CMD ["/prometheus-view-proxy"]