#!/bin/bash

source config.sh

for n in {1..359}
do
   curl ${BACKEND_URL}/api/v1/nationalpark/$n
done
