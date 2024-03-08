# Stage 1: Build Stage
FROM golang:1.21-alpine as build

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/server.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY .env .
COPY --from=build /app/server .
EXPOSE ${PORT}
CMD ["./server"]