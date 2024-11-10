FROM golang:alpine

WORKDIR /app

# Install necessary dependencies
# RUN apk add --no-cache build-base gcc g++

# Copy source code and build
COPY . .
# RUN export CGO_ENABLED=1 && go build -o main .
./main

# Expose the port and set the entrypoint
EXPOSE 8080
CMD ["./main"]
