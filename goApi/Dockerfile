FROM golang:latest

WORKDIR /app

COPY . .

RUN chmod +x entrypoint.sh

RUN cp .env.test .env

RUN go mod download

RUN go build -o main ./cmd/api/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./entrypoint.sh"]