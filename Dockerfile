# Use the official Golang image as the base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Set the entry point for the container
ENTRYPOINT ["./main"]

# Expose the port your application is listening on
EXPOSE 8000