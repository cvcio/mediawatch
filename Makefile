REGISTRY=cvcio
PROJECT=mediawatch-v2
TAG:=$(shell git rev-parse HEAD)
BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
POD=$(shell kubectl get pod -l app=mongo -o jsonpath='{.items[0].metadata.name}')
CONTAINER=$(shell docker ps -f name=mongo -f label=app=mediawatch -q)
BUF_VERSION:=1.28.1
SERVICES=api compare enrich feeds listen scraper worker twitter
NAMESPACE=default

.PHONY: keys
keys: ## generate keys
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

.PHONY: tools
tools: ## install tools
	go get github.com/andefined/realize
	go get github.com/golangci/golangci-lint

.PHONY: buf-install
buf-install: ## install buf version $$BUF_VERSION
	go install github.com/bufbuild/buf/cmd/buf@v${BUF_VERSION}

.PHONY: buf-generate
buf-generate: vendor buf-generate-go buf-generate-tags buf-generate-py buf-clean buf-push ## generate buf files

.PHONY: buf-generate-go
buf-generate-go: ## generate go files
	buf generate --template buf.gen.yaml .

.PHONY: buf-generate-tags
buf-generate-tags: ## generate tagger files
	buf generate --template buf.gen-tags.yaml .

.PHONY: buf-generate-py
buf-generate-py: ## generate python files
	buf generate --template buf.gen.py.yaml \
		--exclude-path proto/mediawatch/articles \
		--exclude-path proto/mediawatch/common \
		--exclude-path proto/mediawatch/compare \
		--exclude-path proto/mediawatch/feeds \
		--exclude-path proto/mediawatch/scrape \
		--exclude-path proto/mediawatch/posts \
		.  

.PHONY: buf-clean
buf-clean: ## clean buf files
	rm -rf pkg/tagger
	rm -rf cmd/enrich/enrich/tagger

.PHONY: buf-update
buf-update: ## update buf modules
	buf mod update

.PHONY: buf-push
buf-push: ## push buf files
	cd proto; buf push

.PHONY: buf-lint
buf-lint: ## lint buf files
	buf lint

.PHONY: vendor
vendor: ## tidy and vendor go modules
	go mod tidy
	go mod vendor

.PHONY: run-mediawatch
run-mediawatch: ## run mediawatch all go services
	realize start

.PHONY: run-api
run-api: ## run api service (deprecated)
	realize start -n api

.PHONY: run-connect-api
run-connect-api: ## run connect-api service (new)
	realize start -n connect-api

.PHONY: run-feeds
run-feeds: ## run feeds service
	realize start -n feeds

.PHONY: run-listen
run-listen: ## run listen service (deprecated)
	realize start -n listen

.PHONY: run-worker
run-worker: ## run worker service
	realize start -n worker

.PHONY: run-compare
run-compare: ## run compare service
	realize start -n compare

.PHONY: run-scraper
run-scraper: ## run scraper service
	cd cmd/scraper; yarn serve

.PHONY: run-twitter
run-twitter: ## run twitter service (deprecated)
	cd cmd/twitter; yarn serve

.PHONY: run-enrich 
run-enrich: ## run enrich service (needs python environment activated)
	cd cmd/enrich; $(MAKE) dev

.PHONY: test-go
test-go: ## run go tests
	go test -v ./...

.PHONY: lint-go
lint-go: ## lint golang code
	golangci-lint run -e vendor -e cmd/scraper -e cmd/enrich

.PHONY: lint-py
lint-py: ## lint python code
	cd cmd/enrich; $(MAKE) lint

.PHONY: lint
lint: lint-go lint-py ## lint all code

.PHONY: db-start
db-start: ## start dbs
	docker compose up -d 

.PHONY: db-logs
db-logs: ## show db logs
	docker compose logs -f

.PHONY: db
db: db-start db-logs ## start dbs and show logs

.PHONY: db-stop
db-stop: ## stop dbs
	docker compose stop

.PHONY: up
up: db-start ## start docker services

.PHONY: down
down: db-stop ## stop docker services

.PHONY: clean-es
clean-es: ## clean elasticsearch
	docker compose stop elasticsearch
	docker compose rm -f -v elasticsearch
	docker volume rm -f mediawatch_data_elasticsearch

.PHONY: clean-kafka
clean-kafka: ## clean kafka
	docker compose stop kafka zookeeper
	docker compose rm -f -v kafka zookeeper
	docker volume rm -f mediawatch_data_kafka mediawatch_data_zookeeper

.PHONY: services-build
services-build: ## build all services
	for name in ${SERVICES}; do\
		cp cmd/$$name/Dockerfile.$$name .;\
		echo "Building image $$name";\
		docker build -f Dockerfile.$$name --rm -t $$name:$(TAG) .;\
		rm Dockerfile.$$name;\
	done

.PHONY: services
services: services-build ## build all services

.PHONY: docker
docker: ## build docker image specified by $$APP
	cp cmd/${APP}/Dockerfile.$(APP) .
	docker build -f Dockerfile.${APP} --rm -t ${APP}:$(TAG) .
	rm Dockerfile.${APP}

.PHONY: prod
prod: vendor ## build docker image specified by $$APP and deploy it to kubernetes
	cp cmd/${APP}/Dockerfile.$(APP) .
	docker build -f Dockerfile.${APP} --rm -t ${APP}:$(TAG) .
	@chmod +x scripts/deploy.sh
	NAME=${APP} REPO=$(REGISTRY) PROJECT=$(PROJECT) NAMESPACE=$(NAMESPACE) CIRCLE_SHA1=$(TAG) CIRCLE_BRANCH=$(BRANCH) scripts/deploy.sh
	rm Dockerfile.${APP}

.PHONY: prod-all
prod-all: vendor ## build all services and deploy them to kubernetes
	for name in ${SERVICES}; do\
		cp cmd/$$name/Dockerfile.$$name .;\
		echo "Building image $$name";\
		docker build -f Dockerfile.$$name --rm -t $$name:$(TAG) .;\
		chmod +x scripts/deploy.sh;\
		NAME=$$name REPO=$(REGISTRY) PROJECT=$(PROJECT) CIRCLE_SHA1=$(TAG) CIRCLE_BRANCH=$(BRANCH) scripts/deploy.sh;\
		rm Dockerfile.$$name;\
	done

.PHONY: get-mongo-backup
get-mongo-backup: ## get mongo backup
	kubectl exec ${POD} -c mongo -i -t -- bash -c 'mongodump -d mediawatch --gzip --archive=/tmp/dump.tar.gz' && kubectl cp ${POD}:/tmp/dump.tar.gz dump.tar.gz

.PHONY: restore-mongo-backup
restore-mongo-backup: ## restore mongo backup
	docker cp dump.tar.gz ${CONTAINER}:/tmp && docker exec ${CONTAINER} mongorestore --drop --gzip --archive=/tmp/dump.tar.gz

# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: release
release: custom 

.PHONY: help
help: ## print help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)