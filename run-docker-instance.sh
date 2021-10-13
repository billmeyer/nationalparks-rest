#!/bin/bash

source config.sh

# See config.sh for setting these values.
docker run -it -p 8080:8080 \
	-e SPLUNK_ACCESS_TOKEN=${SPLUNK_ACCESS_TOKEN} \
	-e SPLUNK_REALM=${SPLUNK_REALM} \
	-e DBHOST=${DBHOST} \
	-e DBPORT=${DBPORT} \
	-e HTTPHOST=${HTTPHOST} \
	-e HTTPPORT=${HTTPPORT} \
	nationalparks-rest

