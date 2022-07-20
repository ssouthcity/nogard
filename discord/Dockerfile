FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o nogard .

FROM alpine

COPY --from=builder /app/nogard /bin

ENTRYPOINT [ "/bin/nogard" ]
