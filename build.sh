#!/bin/bash

source ./server/config.sh

HIGHLIGHT='\033[0;34m'
NC='\033[0m'

echo -e "${HIGHLIGHT}Getting distil-ingest..${NC}"

# get distil-ingest and force a static rebuild of it so that it can run on Alpine
go get -u -v github.com/uncharted-distil/distil-test/
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a github.com/uncharted-distil/distil-test/
mv distil-test ./server

echo -e "${HIGHLIGHT}Building image ${DOCKER_IMAGE_NAME}...${NC}"

# build the docker image
cd server

docker build --squash --no-cache --network=host \
    --tag $DOCKER_REPO/$DOCKER_IMAGE_NAME:${DOCKER_IMAGE_VERSION} \
    --tag $DOCKER_REPO/$DOCKER_IMAGE_NAME:latest .
cd ..


echo -e "${HIGHLIGHT}Done${NC}"
