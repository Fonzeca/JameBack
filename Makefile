# Makefile for building and pushing Docker image

# Variables
IMAGE_NAME = fonzeca/user-hub
TAG = dev

# Default target
all: build push

# Build Docker image
build:
	docker build -t $(IMAGE_NAME):$(TAG) .

# Push Docker image to repository
push:
	docker push $(IMAGE_NAME):$(TAG)

# Clean up local Docker images
clean:
	docker rmi $(IMAGE_NAME):$(TAG)