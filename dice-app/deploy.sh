#!/usr/bin/env bash

set -exu

# docker build -t fagdtw/dice-app .
docker buildx build --platform linux/amd64,linux/arm64 --tag fagdtw/dice-app --push . || echo "Failed to build and push image, asumming one is already published..."
kubectl apply -f ./kubernetes/
