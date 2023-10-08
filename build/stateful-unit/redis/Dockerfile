# Builder stage
ARG BUILDER_IMAGE
FROM $BUILDER_IMAGE as builder

# Setup the Redis image
FROM redis:latest

# Copy the initialization script to the Docker container
COPY init-redis.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/init-redis.sh

COPY --from=builder /app/datagen .
RUN ./datagen redis

# Run Redis and the initialization script during container startup
CMD [ "sh", "-c", "redis-server --daemonize yes --logfile /var/log/redis.log; /usr/local/bin/init-redis.sh; tail -f /var/log/redis.log" ]
