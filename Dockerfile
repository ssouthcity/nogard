FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o nogard cmd/nogard/main.go

FROM alpine

COPY --from=builder /app/nogard /bin

ENTRYPOINT [ "/bin/nogard" ]
