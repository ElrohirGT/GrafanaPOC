#!/usr/bin/env bash

set -exu

docker build -t fagdtw/dice-app .
docker image push fagdtw/dice-app
kubectl apply -f ./kubernetes/
