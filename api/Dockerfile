# Use the official Go image
FROM golang:1.19 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the application
RUN go mod download
RUN go build -o main

# Run the web service on container startup.
CMD ["/app/main"]
