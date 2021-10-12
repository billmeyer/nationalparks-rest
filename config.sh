#!/bin/bash

# IMPORTANT: Uncomment these if they are not set elsewhere in your environment and
# set them to the proper values associated with your Splunk Observability account.
#export SPLUNK_ACCESS_TOKEN=
#export SPLUNK_REALM=

# BACKEND_URL represents the URL where this service will be accessible from.
# It should be in the form (http|https)://(hostname|ip address):(port number).
#export BACKEND_URL=http://localhost:8080
export BACKEND_URL=http://192.168.3.233:80

# The network interface this service should listen on.  Common examples are 0.0.0.0 to listen on all interfaces
# or 127.0.0.1 to only listen on the loopback (non-network accessible) interface.
export HTTPHOST=0.0.0.0

# The network port this service should listen on.  Default is 8080.
export HTTPPORT=8080

# The IP Address the MySQL instance is hosted at.
export DBHOST=192.168.3.230

# The port number the MySQL instance is listening on.
export DBPORT=3306