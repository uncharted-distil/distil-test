#!/bin/bash

source ./server/config.sh
docker push $DOCKER_REPO/$DOCKER_IMAGE_NAME:latest
docker push $DOCKER_REPO/$DOCKER_IMAGE_NAME:${DOCKER_IMAGE_VERSION}
