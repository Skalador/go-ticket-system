# syntax=docker/dockerfile:1
FROM golang:1.21

# Set environment variable
ENV MONGODB_CONNECTION_STRING="mongodb+srv://username:password@database/"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-ticket-system

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8000

# Run
CMD ["/go-ticket-system"]
