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

.PHONY: build push latest deploy rollback

## Сборка образа
build:
	@echo "Building Docker image with tag $(TAG)..."
	docker build --file $(DOCKERFILE) -t $(REGISTRY):$(TAG) .

## Публикация образа
push:
	@echo "Pushing Docker image $(REGISTRY):$(TAG)..."
	docker push $(REGISTRY):$(TAG)

latest:
	@echo "Building and pushing latest image..."
	docker build --file $(DOCKERFILE) -t $(REGISTRY):latest .
	docker push $(REGISTRY):latest

## Деплой на сервер (нужен SSH доступ)
deploy:
	@echo "Deploying version $(TAG) to server..."
	ssh $${VPS_USER}@$${VPS_HOST} "\
		set -e;\
		docker login -u '$${DOCKER_USERNAME}' -p '$${DOCKER_PASSWORD}';\
		cd /path/to/your/app;\
		echo 'Saving current version...';\
		echo \$$(docker ps --filter 'name=myapp' --format '{{.Image}}') > .last_version;\
		docker pull $(REGISTRY):$(TAG);\
		docker compose --file $(DOCKER_COMPOSE_FILE) down;\
		docker compose --file $(DOCKER_COMPOSE_FILE) up -d --remove-orphans;\
	"

## Rollback до предыдущей версии
rollback:
	@echo "Rolling back to previous version..."
	ssh $${VPS_USER}@$${VPS_HOST} "\
		set -e;\
		cd /path/to/your/app;\
		if [ -f .last_version ]; then\
			OLD_IMAGE=$$(cat .last_version);\
			echo 'Rolling back to $$OLD_IMAGE';\
			docker pull $$OLD_IMAGE || true;\
			docker compose --file $(DOCKER_COMPOSE_FILE) down;\
			docker run -d --restart=always --name myapp $$OLD_IMAGE;\
		else\
			echo 'No previous version found for rollback.';\
		fi;\
	"
