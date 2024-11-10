# Build stage
FROM golang:bullseye AS build

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install build-essential gcc g++ -y && export CGO_ENABLED=1 && go build -o main .

# Final stage
FROM alpine

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]