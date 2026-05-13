IMAGE_NAME := "webhook"
IMAGE_TAG := "latest"

build:
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .
