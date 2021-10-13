#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <Docker Hub Username>"
  exit -1
fi

USER=$1
docker tag nationalparks-rest:latest ${USER}/nationalparks-rest:latest
docker push ${USER}/nationalparks-rest:latest
