
FROM golang:1.24.1-alpine as builder
# Set the Current Working Directory inside the container
WORKDIR /go/src/protoconf/demoapp

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /go/bin/main . &&\
    adduser -D -u 10001 demo
# Start a new stage from scratch
FROM scratch AS runtime

# Set the Current Working Directory inside the container
WORKDIR /root/

COPY --from=builder /etc/passwd /etc/passwd

USER 1001
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/main .

COPY --from=builder /go/src/protoconf/demoapp/local/config.json /root/local/config.json

# Expose port 8181 to the outside world
EXPOSE 8181

# Command to run the executable
CMD ["./main"]