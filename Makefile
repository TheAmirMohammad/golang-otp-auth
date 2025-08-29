run: ; go run ./cmd/server
swag: ; go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/server/main.go -o ./docs
tidy: ; go mod tidy