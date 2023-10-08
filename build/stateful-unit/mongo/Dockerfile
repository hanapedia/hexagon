# Builder stage
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

# Final stage
FROM mongo:4.4-focal

WORKDIR /data

COPY --from=builder /app/datagen .
RUN ./datagen mongo

# Copy the custom script and data files
COPY import-data.sh /docker-entrypoint-initdb.d/

# Grant execution permission to the script
RUN chmod +x /docker-entrypoint-initdb.d/import-data.sh

# Default command to run when starting the container
CMD ["mongod"]
