#!/bin/bash

source config.sh

curl --include ${BACKEND_URL}/api/v1/health-check $*
