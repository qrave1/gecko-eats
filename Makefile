# Название образа и путь до registry
REGISTRY=qrave1/gecko-eats
TAG=build_$(shell date '+%Y_%m_%d_%H_%M_%S')

.PHONY: build push generate run latest

# Сборка образа
build:
	@docker build -t $(REGISTRY):$(TAG) .

push:
	@docker push $(REGISTRY):$(TAG)

latest:
	@echo "Building and pushing latest image..."
	@docker build -t $(REGISTRY):latest .
	@docker push $(REGISTRY):latest

generate:
	@wire gen  ./cmd/wire

run:
	@docker compose up --build