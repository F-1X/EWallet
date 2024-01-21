FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY ./ /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /build cmd/app/main.go


FROM alpine

WORKDIR /app

COPY --from=builder /build /build

EXPOSE 3000

ENTRYPOINT ["/build"]

