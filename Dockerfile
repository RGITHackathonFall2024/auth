FROM golang:1.23.2
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
EXPOSE 8080
RUN go build /app/cmd/user-service
CMD ["go", "run", "/app/cmd/user-service"]