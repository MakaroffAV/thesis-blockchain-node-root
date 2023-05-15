# Use the official golang image as base
FROM golang:latest

# Create a directory for the app
WORKDIR /app

# Clone the public repository
RUN git clone https://github.com/MakaroffAV/thesis-blockchain-node-root.git

# Set the working directory to the cloned repository
WORKDIR /app/thesis-blockchain-node-root

# Build the application
RUN go build -o main ./cmd/main.go

# Run the application
CMD [ "./main" ]
