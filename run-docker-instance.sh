#!/bin/bash

source config.sh

docker run -it -p 8080:8080 \
	-e SPLUNK_ACCESS_TOKEN=${SPLUNK_ACCESS_TOKEN} \
	-e SPLUNK_REALM=${SPLUNK_REALM} \
	-e DBHOST=${DBHOST} \
	-e DBPORT=3306 \
	nationalparks-rest

