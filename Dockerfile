FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app 

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/main.go

RUN addgroup -g 10001 -S appgroup && adduser -u 10001 -S -D -G appgroup appuser

FROM alpine:3.19

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/server /server

USER appuser:appgroup

ENTRYPOINT ["./server"]
