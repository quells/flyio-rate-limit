FLY_REGISTRY=registry.fly.io
APP=crimson-tree-7201
TAG=latest
IMG=$(FLY_REGISTRY)/$(APP):$(TAG)

.PHONY: build
build:
	docker build -t $(IMG) .
	docker buildx build --platform linux/amd64 --load -t $(IMG) .

.PHONY: run
run:
	docker run --rm -p 8080:8080 $(IMG)

.PHONY: push
push:
	docker push $(IMG)

.PHONY: deploy
deploy:
	fly deploy -i $(IMG)
