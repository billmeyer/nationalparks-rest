#!/bin/bash

docker run -it -p 8080:8080 \
	-e SPLUNK_ACCESS_TOKEN=${SPLUNK_ACCESS_TOKEN} \
	-e SPLUNK_REALM=${SPLUNK_REALM} \
	-e DBHOST=192.168.3.230 \
	-e DBPORT=3306 \
	nationalparks-rest

