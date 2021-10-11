#!/bin/bash

source config.sh

curl --include -H 'Origin: http://foo.com' ${BACKEND_URL}/api/v1/health-check $*

