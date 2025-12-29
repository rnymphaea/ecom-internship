FROM golang:1.25.5-alpine AS builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/main.go

RUN addgroup -g 10001 -S appgroup && adduser -u 10001 -S -D -G appgroup appuser

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/server /server

USER appuser:appgroup

EXPOSE 8080

ENTRYPOINT ["./server"]
