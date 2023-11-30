FROM golang:1.21-alpine3.17 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o app main.go

FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app /app/app

COPY .env .

EXPOSE 9000

CMD ["/app/app"]