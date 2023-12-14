FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN go build ./cmd/api

FROM scratch

WORKDIR /workspace/app
COPY --from=builder /app/api .
COPY .env .

ENTRYPOINT ["./api"]