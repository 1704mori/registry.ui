FROM golang:1.24.0-alpine3.21 AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o registry_ui ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/registry_ui .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./registry_ui"]