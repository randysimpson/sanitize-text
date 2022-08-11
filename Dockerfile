FROM golang:1.17 as builder
WORKDIR /usr/local/go/src
RUN mkdir sanitize-text
WORKDIR /usr/local/go/src/sanitize-text
ADD . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .
FROM scratch
COPY --from=builder /usr/local/go/src/sanitize-text/main /app/
WORKDIR /app
EXPOSE 8081
CMD ["./main"]
