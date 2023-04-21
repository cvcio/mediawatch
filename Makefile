REGISTRY=cvcio
PROJECT=mediawatch-v2
TAG:=$(shell git rev-parse HEAD)
BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
POD=$(shell kubectl get pod -l app=mongo -o jsonpath='{.items[0].metadata.name}')
CONTAINER=$(shell docker ps -f name=mongo -f label=app=mediawatch -q)
BUF_VERSION:=1.8.0
SERVICES=api compare enrich feeds listen scraper worker
NAMESPACE=default

keys:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

tools:
	go get github.com/oxequa/realize
	go get github.com/golangci/golangci-lint

buf-install:
	curl -sSL \
    	"https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(shell uname -s)-$(shell uname -m)" \
    	-o "$(shell go env GOPATH)/bin/buf" && \
  	chmod +x "$(shell go env GOPATH)/bin/buf"

buf-generate: vendor buf-generate-go buf-generate-tags buf-generate-py buf-clean

buf-generate-go:
	buf generate --template buf.gen.yaml .

buf-generate-tags:
	buf generate --template buf.gen-tags.yaml .

buf-generate-py:
	buf generate --template buf.gen.py.yaml \
		--exclude-path proto/mediawatch/articles \
		--exclude-path proto/mediawatch/common \
		--exclude-path proto/mediawatch/compare \
		--exclude-path proto/mediawatch/feeds \
		--exclude-path proto/mediawatch/scrape \
		.  
buf-clean:
	rm -rf pkg/tagger
	rm -rf cmd/enrich/enrich/tagger

buf-update:
	buf mod update

buf-lint:
	buf lint

vendor:
	go mod vendor

run-mediawatch:
	realize start

run-api:
	realize start -n api

run-connect-api:
	realize start -n connect-api

run-feeds:
	realize start -n feeds

run-listen:
	realize start -n listen

run-worker:
	realize start -n worker

run-compare:
	realize start -n compare

run-ministry:
	realize start -n ministry

run-scraper:
	cd cmd/scraper; yarn serve

run-enrich:
	cd cmd/enrich; $(MAKE) dev

test:
	go test -v ./...

lint:
	golangci-lint run -e vendor

db-start:
	docker-compose up -d

db-logs:
	docker-compose logs -f

db: db-start db-logs

db-stop:
	docker-compose stop

clean-es:
	docker-compose stop elasticsearch
	docker-compose rm -f -v elasticsearch
	docker volume rm -f mediawatch_data_elasticsearch

services-build:
	for name in ${SERVICES}; do\
		cp cmd/$$name/Dockerfile.$$name .;\
		echo "Building image $$name";\
		docker build -f Dockerfile.$$name --rm -t $$name:$(TAG) .;\
		rm Dockerfile.$$name;\
	done

services-run:
	docker-compose -f docker-compose.with-services.yaml up

services: services-build services-run

prod:
	go mod vendor
	cp cmd/${APP}/Dockerfile.$(APP) .
	docker build -f Dockerfile.${APP} --rm -t ${APP}:$(TAG) .
	@chmod +x scripts/deploy.sh
	NAME=${APP} REPO=$(REGISTRY) PROJECT=$(PROJECT) NAMESPACE=$(NAMESPACE) CIRCLE_SHA1=$(TAG) CIRCLE_BRANCH=$(BRANCH) scripts/deploy.sh
	rm Dockerfile.${APP}

prod-all:
	go mod vendor
	for name in ${SERVICES}; do\
		cp cmd/$$name/Dockerfile.$$name .;\
		echo "Building image $$name";\
		docker build -f Dockerfile.$$name --rm -t $$name:$(TAG) .;\
		chmod +x scripts/deploy.sh;\
		NAME=$$name REPO=$(REGISTRY) PROJECT=$(PROJECT) CIRCLE_SHA1=$(TAG) CIRCLE_BRANCH=$(BRANCH) scripts/deploy.sh;\
		rm Dockerfile.$$name;\
	done

get-mongo-backup:
	kubectl exec ${POD} -c mongo -i -t -- bash -c 'mongodump -d mediawatch --gzip --archive=/tmp/dump.tar.gz' && kubectl cp ${POD}:/tmp/dump.tar.gz dump.tar.gz

restore-mongo-backup:
	docker cp dump.tar.gz ${CONTAINER}:/tmp && docker exec ${CONTAINER} mongorestore --drop --gzip --archive=/tmp/dump.tar.gz

# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: release
release: custom 
