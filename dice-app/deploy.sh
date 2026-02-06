#!/usr/bin/env bash

set -exu

docker build -t fagdtw/dice-app .
docker image push fagdtw/dice-app | echo "Failed to publish the image, assuming someone published it before you..."
kubectl apply -f ./kubernetes/
