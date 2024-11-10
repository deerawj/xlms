FROM anjanamadu/go1.21-gcc-alpine:latest

WORKDIR /app

# Copy source code and build
COPY . .

# Expose the port and set the entrypoint
EXPOSE 8080
CMD ["/bin/sh", "./entrypoint.sh"]
