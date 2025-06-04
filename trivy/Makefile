.PHONY: scan-image

scan-with-docker:
	@if [ -n "$(image_name)" ]; then \
		IMAGE_NAME="$$(echo $(image_name) | cut -d ':' -f1)"; \
		IMAGE_TAG="$$(echo $(image_name) | cut -d ':' -f2)"; \
	elif [ -f .env ]; then \
		IMAGE_NAME="$$(grep IMAGE_NAME .env | cut -d '=' -f2 | tr -d '[:space:]')"; \
		IMAGE_TAG="$$(grep IMAGE_TAG .env | cut -d '=' -f2 | tr -d '[:space:]')"; \
	else \
		echo "Error: No image_name provided and .env file not found"; \
		exit 1; \
	fi; \
	echo "Scanning image $$IMAGE_NAME:$$IMAGE_TAG..."; \
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy:0.63.0 image "$$IMAGE_NAME:$$IMAGE_TAG"