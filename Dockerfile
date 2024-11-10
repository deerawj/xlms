# Build stage
FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN export CGO_ENABLED=1 && go build -o main .

# Final stage
FROM alpine

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]