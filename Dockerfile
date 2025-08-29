FROM golang:1.25 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/server/main.go -o ./docs
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/otp-service ./cmd/server

FROM gcr.io/distroless/base-debian12
COPY --from=build /bin/otp-service /usr/local/bin/otp-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/otp-service"]
