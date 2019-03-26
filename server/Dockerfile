# Set the base image to alpine
FROM alpine:3.9.2

# make sure package db is up to date
RUN apk update

# missing from baseline alpine
RUN apk add --no-cache ca-certificates

# download build toolchain
RUN apk add --no-cache --virtual .build-deps \
  gcc \
  linux-headers \
  make \
  musl-dev \
  tar

# copy in distil-ingest
RUN mkdir -p /distil
WORKDIR /distil
ADD config.sh .
ADD distil-test .

# copy run script across
ADD start_test.sh /

CMD . /start_test.sh
