.PHONY: up down scan-cli

up:
	docker compose up -d

down:
	docker compose down

scan-cli:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@if [ -z "$(input_dir)" ]; then \
		echo "Error: input_dir not specified"; \
		echo "Usage: make scan input_dir=/path/to/code"; \
		exit 1; \
	fi
	SONAR_URL=$$(grep SONARQUBE_URL .env | cut -d '=' -f2); \
	SONAR_TOKEN=$$(grep SONAR_TOKEN .env | cut -d '=' -f2); \
	echo "Starting SonarQube scan..."; \
	echo "Scanning $(input_dir) with SonarQube at $$SONAR_URL"; \
	docker run --rm \
		-v "$(input_dir):/usr/src" \
		-e SONAR_HOST_URL="$$SONAR_URL" \
		-e SONAR_TOKEN="$$SONAR_TOKEN" \
		sonarsource/sonar-scanner-cli