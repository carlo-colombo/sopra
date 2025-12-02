# Use the official Golang image as a base image
FROM golang:1.22-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
# The CGO_ENABLED=0 flag disables cgo which is necessary for static linking
# The -o sopra flag sets the output file name to "sopra"
RUN go build -o sopra ./main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/app/sopra"]
