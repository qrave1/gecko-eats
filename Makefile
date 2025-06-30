# Название образа и путь до registry
REGISTRY=qrave1/gecko-eats
TAG=build_$(shell date '+%Y_%m_%d_%H_%M_%S')

.PHONY: build push generate

# Сборка образа
build:
	docker build -t $(REGISTRY):$(TAG) .

push:
	docker push $(REGISTRY):$(TAG)

generate:
	@wire gen  ./cmd/wire

run:
	@docker compose up --build