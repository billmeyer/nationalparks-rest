#!/bin/bash

for n in {1..360}
do
   curl http://localhost:8080/api/v1/nationalpark/$n
done