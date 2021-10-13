#!/bin/bash

version=$(date +%s)
docker build -t nationalparks-rest:${version} . $*
docker tag nationalparks-rest:${version} nationalparks-rest:latest
