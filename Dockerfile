FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o gwi-server ./cmd/gwiapp
WORKDIR /app

RUN cp /build/gwi-server .

FROM scratch

COPY --from=builder /app/gwi-server /

EXPOSE 31000

ENTRYPOINT ["/gwi-server"]