#!/bin/bash

source ./server/config.sh

docker run \
  --rm \
  --name $DOCKER_IMAGE_NAME \
  $DOCKER_REPO/$DOCKER_IMAGE_NAME:${DOCKER_IMAGE_VERSION}
