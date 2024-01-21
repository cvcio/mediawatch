REGISTRY=cvcio
PROJECT=mediawatch-svc
TAG=`cat VERSION`
PROTO_VERSION=1.2.0
MODULE=enrich

REG_PROJ=$(REGISTRY)/$(PROJECT)-$(MODULE)
REG_TAG=$(REGISTRY)/$(PROJECT)-$(MODULE):$(TAG)

CURRENT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUF_VERSION=1.1.0

help: ## Print Help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start server in dev mode
	nodemon -e py -x python enrich/app.py

buf-install: ## Install buf
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
  	chmod +x "$(shell go env GOPATH)/bin/buf"

proto-get: ## Get protocol buffers
	curl -sSL \
    	https://github.com/plagiari-sm/svc-proto/archive/$(PROTO_VERSION).tar.gz | tar xz && \
    	mv svc-proto-$(PROTO_VERSION) proto

proto-build: ## Generate protocol buffers
	buf generate

proto-cleanup: ## Remove protocol buffers
	rm -rf proto

proto: proto-get proto-build proto-cleanup ## Get & Build protocol buffers

docker: # Build docker image
	docker build -f Dockerfile --rm -t $(REG_TAG) .
	docker tag $(REG_TAG) $(REG_PROJ):latest

docker-push: # Push docker image
	docker push $(REG_TAG)
	docker push $(REG_PROJ):latest

deployment: # Update the deployment
	kubectl -n default set image deployment/${MODULE} ${MODULE}=$(REG_TAG)

prod: docker docker-push deployment ## Publish to production

lint: ## Lint the code
	PYTHONPATH=. pylint --fail-under=9 enrich

.DEFAULT_GOAL := help
.PHONY: help dev buf-install proto proto-get proto-build proto-cleanup docker docker-push deployment prod
