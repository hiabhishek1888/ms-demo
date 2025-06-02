FROM golang:1.24.3-alpine

WORKDIR /app

# Create the go.mod file
COPY go.mod ./

RUN go get github.com/gin-gonic/gin
RUN go get github.com/lib/pq

# Copy the source code
COPY main.go ./

# Download all dependencies and verify
RUN go mod tidy

# Build the application
RUN go build -o main .

# Expose port 8081
EXPOSE 8081

# Command to run the application
CMD ["./main"] 