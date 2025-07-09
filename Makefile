# Название образа и путь до registry
DOCKERFILE=deploy/Dockerfile
DOCKER_COMPOSE_FILE=deploy/docker-compose.yml
REGISTRY=qrave1/gecko-eats
TAG=build_$(shell date '+%Y_%m_%d_%H_%M_%S')

.PHONY: build push run latest

# Сборка образа
build:
	@docker build --file $(DOCKERFILE) -t $(REGISTRY):$(TAG) .

push:
	@docker push $(REGISTRY):$(TAG)

latest:
	@echo "Building and pushing latest image..."
	@docker build --file $(DOCKERFILE) -t $(REGISTRY):latest .
	@docker push $(REGISTRY):latest

run_infra:
	@docker compose --file $(DOCKER_COMPOSE_FILE) --profile infra up -d --remove-orphans
