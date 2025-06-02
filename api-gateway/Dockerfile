FROM golang:1.24.3-alpine

WORKDIR /app

# Create the go.mod file
COPY go.mod ./

RUN go get github.com/gin-gonic/gin

# Copy the source code
COPY main.go ./

# Download all dependencies and verify
RUN go mod tidy

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 