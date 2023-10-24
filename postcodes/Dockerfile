FROM golang:1.21-alpine AS builder
#RUN apk --no-cache add ca-certificates
RUN apk update && apk add --no-cache git

# Set necessary environmet variables needed for our image
#ENV GOOS=linux \
#    GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./

# Installs Go dependencies
RUN go mod download

# Copy the entire project and build it
# This layer is rebuilt when a file changes in the project directory
# We want to populate the module cache based on the go.{mod,sum} files.
COPY . .
COPY .env /.env


# Builds your app with optional configuration
RUN go build -v -o postcodes

# This results in a single layer image
FROM scratch

# Export necessary port
EXPOSE 8089

COPY --from=builder /app/postcodes /postcodes
ENTRYPOINT ["/postcodes"]
CMD ["--help" ]

