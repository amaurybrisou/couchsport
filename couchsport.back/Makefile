ifndef GOROOT 
$(error "GOROOT is not set")
endif

ifndef GOBIN
$(error "GOBIN is not set")
endif

ENV=local

RELEASE_PATH=$(CURDIR)/release

GOCMD=$(GOROOT)/bin/go
DOCKERCMD=$(shell which docker)

DOCKER_DEV_IMAGE=brisouamaury/couchsport.back
DOCKER_PROD_IMAGE=brisouamaury/couchsport.back.prod

DOCKERFILE_PROD=Dockerfile
DOCKERFILE_DEV=Dockerfile


.DEFAULT_GOAL := run-soft

## build: build production binary in ${RELEASE_PATH}
build:
	[ -d $(RELEASE_PATH) ] || mkdir $(RELEASE_PATH)
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build --no-cache . -f $(DOCKERFILE_PROD) --target compile --output $(RELEASE_PATH) --platform $(ENV) --progress plain

## build-prod: build production container for linux
build-prod:
	[ -d $(RELEASE_PATH) ] || mkdir $(RELEASE_PATH)
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build --no-cache . -f $(DOCKERFILE_PROD) --target release -t $(DOCKER_PROD_IMAGE) --platform linux --progress plain

## run: Run production container
run: build-prod
	DOCKER_BUILDKIT=1 $(DOCKERCMD) run --env-file .env -t $(DOCKER_PROD_IMAGE) $(DOCKER_PROD_IMAGE)

## stop-prod: Stop production container
stop-prod:
	$(DOCKERCMD) stop $(DOCKER_PROD_IMAGE) & $(DOCKERCMD) rm $(DOCKER_PROD_IMAGE)



## build-dev: Build development container
build-dev:
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build . -f $(DOCKERFILE_DEV) --target run-dev -t $(DOCKER_DEV_IMAGE) --progress plain

run-soft:
	$(DOCKERCMD) run -it -p 9000:9000 -v ${PWD}/:/go/app --env-file .env --name $(DOCKER_DEV_IMAGE) $(DOCKER_DEV_IMAGE)

## run-dev: Run dev container
run-dev: build-dev run-soft

## stop-dev: Stop development container
stop-dev:
	$(DOCKERCMD) stop $(DOCKER_DEV_IMAGE) && $(DOCKERCMD) rm $(DOCKER_DEV_IMAGE)

## clean: Clean environment
clean:
	[ -d $(RELEASE_PATH) ] && rm -rf $(RELEASE_PATH)

help: Makefile
	@echo
	@echo "Available commands:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo