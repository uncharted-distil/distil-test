#!/bin/bash

source ./server/config.sh
docker login $DOCKER_REPO
docker push $DOCKER_REPO/$DOCKER_IMAGE_NAME:latest
docker push $DOCKER_REPO/$DOCKER_IMAGE_NAME:${DOCKER_IMAGE_VERSION}
docker logout $DOCKER_REPO
