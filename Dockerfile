# Build stage
FROM golang:1.20-alpine3.17 AS builder

# define working directory
WORKDIR /app
# copy whole file in repo to current workdir
COPY . .

# download curl
RUN apk add curl
# build executional file
RUN go build -o myapp main.go
# download migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.17

WORKDIR /app
# copy binary file from builder to current workdir
COPY --from=builder /app/myapp .
# copy migrate
COPY --from=builder /app/migrate .
COPY db/migration ./migration
COPY start.sh .
COPY app.env .

EXPOSE 3000

CMD ["/app/myapp"]
ENTRYPOINT ["/app/start.sh"]