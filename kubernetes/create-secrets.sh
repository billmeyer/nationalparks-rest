#!/bin/bash

source ../config.sh

# Creates a secret named `splunk-access` in the `nationalparks` namespace that will be used to store the
# secrets needed to access Splunk Observability.

# Creates two keys:
#
# realm: Stores the Splunk Observability that will be used
# token: Stores the Acess Token used to access Splunk Obserability
#
# These values are read from environment variables.

if [[ -z "${SPLUNK_ACCESS_TOKEN}" ]]; then
  echo "Environment variable SPLUNK_ACCESS_TOKEN not defined."
  exit -1
fi

if [[ -z "${SPLUNK_REALM}" ]]; then
  echo "Environment variable SPLUNK_REALM not defined."
  exit -1
fi

kubectl create secret generic --namespace nationalparks splunk-access \
  --from-literal=token="${SPLUNK_ACCESS_TOKEN}" \
  --from-literal=realm="${SPLUNK_REALM}"
