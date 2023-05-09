REGISTRY=cvcio
PROJECT=mediawatch-svc
TAG=`cat VERSION`
MODULE=twitter

REG_PROJ=$(REGISTRY)/$(PROJECT)-$(MODULE)
REG_TAG=$(REGISTRY)/$(PROJECT)-$(MODULE):$(TAG)

CURRENT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

docker:
	docker build -f Dockerfile --rm -t $(REG_TAG) .
	docker tag $(REG_TAG) $(REG_PROJ):latest

docker-push:
	docker push $(REG_TAG)
	docker push $(REG_PROJ):latest

deployment:
	kubectl -n default set image deployment/${MODULE} ${MODULE}=$(REG_TAG)

prod: docker docker-push deployment
