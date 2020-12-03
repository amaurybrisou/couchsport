ifndef GOROOT 
$(error "GOROOT is not set")
endif

ifndef GOBIN
$(error "GOBIN is not set")
endif

ENV=local
PROJECTNAME=$(shell basename "$(PWD)")
RELEASE_PATH=$(CURDIR)/release

GOCMD=$(GOROOT)/bin/go
DOCKERCMD=$(shell which docker)

DOCKER_IMAGE_DEV=couchsport.front
DOCKER_IMAGE_PROD=couchsport.front.prod

## build development docker image
build-dev:
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build . --platform $(ENV) --progress plain -t $(DOCKER_IMAGE_DEV)

## run development docker image
run-dev:
	docker run -it -v ${PWD}:/usr/src/app -v /usr/src/app/node_modules -p 8081:8081 $(PROJECTNAME)

## build production docker image
build: build-prod
build-prod:
	DOCKER_BUILDKIT=1 $(DOCKERCMD) build . --target release --progress plain -t $(DOCKER_IMAGE_PROD)


help: Makefile
	@echo
	@echo "Available commands:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo