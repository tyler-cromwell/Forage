FROM golang:1.24.0-alpine3.21

# Install project-specific code
WORKDIR /root
COPY . .

# Run setup commands
RUN go build
RUN source .env

EXPOSE 8001

# Default command
CMD ["./forage"]