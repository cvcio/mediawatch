#!/bin/bash

echo "Running '$0' for '${NAME}'"

### The following section will push your container image to ECR. The `$NAME` variable is provided from our
### Makefile under 'deploy:' rule, which is set to the name of the component/module/service.
###
docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:${CIRCLE_BRANCH}
#docker login -u="$DOCKER_USER" -p="$DOCKER_PASS" ${REPO}
docker push ${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
docker push ${REPO}/${PROJECT}-${NAME}:${CIRCLE_BRANCH}
if [ "$CIRCLE_BRANCH" = "master" ]
  then
  docker tag ${NAME}:${CIRCLE_SHA1} ${REPO}/${PROJECT}-${NAME}:latest
  docker push ${REPO}/${PROJECT}-${NAME}:latest
fi

if [ "${CIRCLE_BRANCH}" = "master" ]
  then
  kubectl -n default set image deployment/${NAME} ${NAME}=${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
fi

if [ "${CIRCLE_BRANCH}" = "development" ]
  then
  kubectl -n development set image deployment/${NAME} ${NAME}=${REPO}/${PROJECT}-${NAME}:${CIRCLE_SHA1}
fi
