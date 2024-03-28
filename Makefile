# Project Vars
APP_NAME=jeopardy-data-scraper
DOCKER_IMAGE_NAME=georgedinicola/$(APP_NAME)
DOCKER_TAG=latest

BUILD_PATH=./cmd/$(APP_NAME)
DB_TABLE_NAME=jeopardy_game_box_scores
CHIP_ARCH=arm64
OS=linux

# Env Vars
DB_HOST ?= localhost
DB_USERNAME ?= default_user
DB_PASSWORD ?= default_password
DB_NAME ?= default_db


# Execution
run: build
	docker run \
		-e APP_MODE="INCREMENTAL" \
		-e DB_HOST=$(DB_HOST) \
		-e DB_USERNAME=$(DB_USERNAME) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(DB_NAME) \
		$(DOCKER_IMAGE_NAME)

run-full: build
	docker run \
		-e APP_MODE="FULL" \
		-e DB_HOST=$(DB_HOST) \
		-e DB_USERNAME=$(DB_USERNAME) \
		-e DB_PASSWORD=$(DB_PASSWORD) \
		-e DB_NAME=$(DB_NAME) \
		$(DOCKER_IMAGE_NAME)

run-excel: build
	docker run \
		-e APP_MODE="EXCEL" \
		-v ~/Desktop:/data \
		$(DOCKER_IMAGE_NAME)

stop:
	docker stop $(DOCKER_IMAGE_NAME)


# Testing
test:
	go test ./... -cover -v

query-postgres-test:
	docker exec -it postgres psql -U $(DB_USERNAME) -d $(DB_NAME)  -c "SELECT * FROM $(TABLE_NAME)"


# Build
deps:
	go mod tidy && \
    go mod verify

.PHONY: build
build:
	docker build --platform $(OS)/$(CHIP_ARCH) --tag $(DOCKER_IMAGE_NAME) .

build-and-run: build run

docker-tag:
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

go-build:
	go build -o $(APP_NAME) GOOS=$(OS) GOARCH=$(CHIP_ARCH) -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go


# Infrastructure
start-postgres:
	docker run -d --name postgres -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p 5432:5432 postgres

stop-postgres:
	docker stop postgres

remove-postgres: stop-postgres
	docker rm postgres
tear-down-postgres: remove-postgres


# Utilities
security-scan:
	gosec ./...

check-dependencies:
	@which docker > /dev/null || (echo "docker is not installed" && exit 1) || \
	@which go > /dev/null || (echo "Go is not installed" && exit 1)