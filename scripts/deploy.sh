#!/bin/bash

echo "Running '$0' for '${NAME}'"

### The following section will push your container image to ECR. The `$NAME` variable is provided from our
### Makefile under 'deploy:' rule, which is set to the name of the component/module/service.
###
docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:${CIRCLE_BRANCH}

docker push ${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
docker push ${REPO}/${PROJECT}-${NAME}:${CIRCLE_BRANCH}

if [ "$CIRCLE_BRANCH" = "main" ]
    then
    docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:latest
    docker push ${REPO}/${PROJECT}-${NAME}:latest
fi

# docker images
# `aws ecr get-login --no-include-email --region ap-northeast-1`
# docker push 963826138034.dkr.ecr.ap-northeast-1.amazonaws.com/${NAME}:${CIRCLE_SHA1}

### If you need to deploy your service to kubernetes, you need to use the kubectl tool. The setup for kubectl's
### config file is also included in this section, although you need to provide your own kubectl's config file
### in the root directory. In this example, the filename is `kubeconf.yaml`.
###
#curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
#chmod +x ./kubectl
#mv ./kubectl /usr/local/bin/kubectl
#mkdir ${HOME}/.kube
## cp kubeconf.yaml ${HOME}/.kube/config
#echo $KUBE_CONFIG >> ${HOME}/.kube/config

### Deploy
#if [ "${CIRCLE_BRANCH}" = "master" ]
#  then
#  kubectl -n ${NAMESPACE} set image deployment/${NAME} ${NAME}=${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
#fi
#
#if [ "${CIRCLE_BRANCH}" = "development" ]
#  then
#  kubectl -n development set image deployment/${NAME} ${NAME}=${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
#fi
