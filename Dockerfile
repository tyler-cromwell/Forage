FROM golang:1.24.0-alpine3.21

# Install project dependencies
WORKDIR /root
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run setup commands
RUN go build -o forage .
RUN source .env

EXPOSE 8001

# Default command
CMD ["./forage"]