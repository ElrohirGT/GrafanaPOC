#!/usr/bin/env bash

set -exu

helm repo add minio https://charts.min.io/
helm upgrade --install minio --namespace monitoring minio/minio --values ./values.yml --wait

brew install minio/stable/mc || echo "Installing the MinIO client..."
sleep 3s
kubectl -n monitoring port-forward svc/minio 9000:9000 &
sleep 1s
BG_PROC=$!

mc alias set minio-local http://localhost:9000 admin password123
mc mb minio-local/loki-chunks
mc mb minio-local/loki-ruler
mc mb minio-local/loki-admin
mc mb minio-local/tempo-traces

echo "Checkout buckets..."
mc ls minio-local

kill "$BG_PROC"
