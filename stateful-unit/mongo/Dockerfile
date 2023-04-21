# Builder stage
FROM golang:latest as builder

WORKDIR /app

# Copy the Go script
COPY generate-dummy-data.go .

# Run the Go script
RUN go run generate-dummy-data.go

# Final stage
FROM mongo:4.4-focal

# Copy the custom script and data files
COPY import-data.sh /docker-entrypoint-initdb.d/
COPY --from=builder /app/*.json /data/

# Grant execution permission to the script
RUN chmod +x /docker-entrypoint-initdb.d/import-data.sh

# Default command to run when starting the container
CMD ["mongod"]
