run: ; go run ./cmd/server
swag: ; go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/server/main.go -o ./docs
tidy: ; go mod tidy
docker: ; docker build -t otp-service:dev .
compose-build: ; docker compose build --no-cache

# full reset to avoid stale containers
## Note: `up` runs in detached mode
up:
	docker compose down --remove-orphans
	docker compose up -d --build
## Note: `up-attached` runs in attached mode
up-attached:
	# full reset to avoid stale containers
	docker compose down --remove-orphans
	docker compose up -d --build
## Note: `up-api` only attaches api service, useful for debugging api only
up-api:
	# full reset to avoid stale containers
	docker compose down --remove-orphans
	docker compose up --build --no-attach db --no-attach redis
## Note: `up-no-build` uses existing images, useful for quick restarts
up-no-build:
	# full reset to avoid stale containers
	docker compose down --remove-orphans
	docker compose up -d

down: ; docker compose down --remove-orphans
down-clean: docker compose down -v --remove-orphans