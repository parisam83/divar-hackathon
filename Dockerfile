FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Set environment variables
ENV GO111MODULE=on

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o subway-finder

# Set working directory in the new container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/subway-finder .

EXPOSE 8080

# Run the application
CMD ["./subway-finder"]
