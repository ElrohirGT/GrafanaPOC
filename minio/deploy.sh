#!/usr/bin/env bash

set -exu

helm repo add minio https://charts.min.io/
helm install minio --namespace monitoring minio/minio --values ./values.yml
