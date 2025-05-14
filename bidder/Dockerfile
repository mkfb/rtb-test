
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY ./bidder .
RUN go mod init rtb
RUN go mod download
RUN go build -o rtb

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rtb /app/rtb
EXPOSE 8080

CMD ["/app/rtb"]