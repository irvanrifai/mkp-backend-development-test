# Variabel
BINARY_NAME=go-boilerplate
MAIN_PATH=cmd/api/main.go

## All: Menjalankan urutan standar
all: deps build run

## Deps: Download dependencies & merapikan go.mod
deps:
	@echo "Tidying and downloading dependencies..."
	go mod tidy
	go mod download

## Build: Mengompilasi aplikasi ke folder tmp
build:
	@echo "Building binary..."
	mkdir -p tmp
	go build -o tmp/$(BINARY_NAME) $(MAIN_PATH)

## Run: Menjalankan aplikasi secara langsung
run:
	@echo "Running application..."
	go run $(MAIN_PATH)

## Test: Run all unit tests
test:
	go test -v ./...

## Test-Cover: Run tests with coverage report
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Gen-Mock: Generate mock menggunakan mockery
gen-mock:
	@echo "Generating mock files..."
	go run github.com/vektra/mockery/v2@latest --name=IUserRepository --dir=internal/repository/user --output=mocks --outpkg=mocks --case=underscore --with-expecter
	go run github.com/vektra/mockery/v2@latest --name=IUserUsecase --dir=internal/usecase/user --output=mocks --outpkg=mocks --case=underscore --with-expecter
	go run github.com/vektra/mockery/v2@latest --name=IScheduleRepository --dir=internal/repository/schedule --output=mocks --outpkg=mocks --case=underscore --with-expecter
	go run github.com/vektra/mockery/v2@latest --name=IScheduleUsecase --dir=internal/usecase/schedule --output=mocks --outpkg=mocks --case=underscore --with-expecter

## Clean: Menghapus file binary
clean:
	@echo "Cleaning up..."
	rm -rf tmp/

.PHONY: all deps build run test gen-mock clean

include .env
export

# URL Database
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

## Migration: Create new migration with TIMESTAMP
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir database/migrations -format "20060102150405" $${name}

## Migration: Run all up migrations
migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

## Migration: Rollback 1 migration
migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down 1

.PHONY: migrate-create migrate-up migrate-down
