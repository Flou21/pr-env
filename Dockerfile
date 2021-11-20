FROM golang:1.16.6-alpine3.14 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build .

CMD /app/pr-env

FROM alpine:3.14.0

WORKDIR /app

COPY --from=builder /app/pr-env pr-env

CMD /app/pr-env