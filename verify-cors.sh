#!/bin/bash

#HOST=localhost:8080
HOST=192.168.3.233:80  # <=== replace with the address of the Kubernetes deployment

curl --include -H 'Origin: http://foo.com' http://${HOST}/api/v1/health-check $*

