# Конфигурация
DOCKERFILE=deploy/Dockerfile
DOCKER_COMPOSE_FILE=deploy/docker-compose.yml
REGISTRY=qrave1/gecko-eats
TAG=$(shell git rev-parse --short HEAD)

# Переменные среды (пригодятся в CI)
export DOCKERFILE
export DOCKER_COMPOSE_FILE
export REGISTRY
export TAG

.PHONY: build push

## Создание миграции
migrate-create:
	@goose -dir internal/infrastructure/postgres/migrations create 1 sql

## Сборка образа
build:
	@echo "Building Docker image with tag $(TAG)..."
	docker build --file $(DOCKERFILE) -t $(REGISTRY):$(TAG) .

## Публикация образа
push:
	@echo "Pushing Docker image $(REGISTRY):$(TAG)"
	docker build --file $(DOCKERFILE) -t $(REGISTRY):$(TAG) .
	docker push $(REGISTRY):$(TAG)
