DOCKER_COMPOSE_FILE_PROD=docker-compose-prod.yml
DOCKER_COMPOSE_FILE_DEV=docker-compose.yml

ENV_FILE_DEV=.env.development
ENV_FILE_PROD=.env.production

build-dev: build-back-dev build-front-dev
build-prod: build-back-prod build-front-prod

## build-front-dev: build the frontend docker image for dev purpose
build-front-prod:
	@echo "Building Front"
	cd couchsport.front/ && make build

## build-front-dev: build the frontend docker image for dev purpose
build-front-dev:
	@echo "Building Front"
	cd couchsport.front/ && make build-dev

## start-dev: Start in development mode. Gets reloaded automatically when code changes.
start-dev:
	@echo "Running Servers..."
	docker network create external-net || true
	docker-compose --env-file $(ENV_FILE_DEV) -f $(DOCKER_COMPOSE_FILE_DEV) up -d

## start: in production
## start-prod: in production
start: start-prod
start-prod:
	@echo "Starting in production..."
	docker network create external-net || true
	docker-compose --env-file $(ENV_FILE_PROD) -f $(DOCKER_COMPOSE_FILE_PROD) up -d

## build-back-prod: build the backend docker image for prod purpose
build-back-prod:
	@echo "Building Go Binary..."
	cd couchsport.back/ && make build-prod

## build-back-dev: build the backend docker image for dev purpose
build-back-dev:
	@echo "Building Go Binary..."
	cd couchsport.back/ && make build-dev

## stop: Stop all dev dockers.
stop-dev:
	@echo "Stopping Servers..."
	docker-compose --env-file $(ENV_FILE_DEV) down

## stop: Stop all dockers.
stop:
	@echo "Stopping Servers..."
	docker-compose --env-file $(ENV_FILE_PROD) down

arg = $(filter-out $@,$(MAKECMDGOALS))
## restart: Restart the specified unit
restart:
	@if [ ! -z "$(arg)" ]; then echo "Restarting unit $(arg)" && docker-compose --env-file $(ENV_FILE_DEV) restart $(arg);	else echo "Service name required as first argument" && docker-compose ps; fi

## logs: Display logs in the console.
logs:
	docker-compose --env-file $(ENV_FILE_DEV) logs -f

## clean: Remove all unused locale Volumes and remove all stopped containers.
clean:
	docker system prune -f
	docker volume prune -f

help: Makefile
	@echo
	@echo "Available commands:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo