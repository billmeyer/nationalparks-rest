#!/bin/bash

#HOST=localhost:8080
HOST=192.168.3.233:80  # <=== replace with the address of the Kubernetes deployment

for n in {1..359}
do
   curl http://${HOST}/api/v1/nationalpark/$n
done