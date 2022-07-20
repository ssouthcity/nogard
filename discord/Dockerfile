FROM golang:1.18 AS builder

WORKDIR /app

COPY . .

RUN go build -o nogard .

FROM alpine

COPY --from=builder /app/nogard /bin

ENTRYPOINT [ "/bin/nogard" ]
