FROM golang:1.16.7-buster AS builder

WORKDIR /root

COPY go.mod go.sum /root/
RUN go mod download

COPY server.go /root/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -extldflags '-static'" server.go

FROM scratch AS runner

COPY --from=builder /root/server /
CMD ["/server"]
