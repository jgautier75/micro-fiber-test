FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /app

#Copy go modules and go sum
COPY go.mod ./
COPY go.sum ./

#Download dependencies
RUN go mod download

# Copy sources
COPY *.go ./
COPY pkg ./pkg

# Build module
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o micro-fiber-test

FROM scratch

WORKDIR /app

COPY --from=builder /app/micro-fiber-test /usr/bin/

EXPOSE 8443

# Run the server executable
CMD [ "/app/micro-fiber-test" ]
