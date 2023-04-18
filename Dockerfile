# Build stage
FROM golang:1.20-alpine3.17 AS builder

# define working directory
WORKDIR /app
# copy whole file in repo to current workdir
COPY . .

# build executional file
RUN go build -o myapp main.go

# Run stage
FROM alpine:3.17

WORKDIR /app
# copy binary file from builder to current workdir
COPY --from=builder /app/myapp .
COPY app.env .

EXPOSE 3000

CMD ["/app/myapp"]