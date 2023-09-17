# Use the official Go image
FROM golang:1.19 as builder

RUN apt-get update && apt-get install -y git

# Copy local code to the container image.
COPY . /app
WORKDIR /app

# Build the application
RUN go mod download
RUN go build -o main .

# Run the web service on container startup.
CMD ["/app/startup.sh"]
