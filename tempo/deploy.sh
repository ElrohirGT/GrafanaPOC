#!/usr/bin/env bash

set -exu

kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.1/cert-manager.yaml
sleep 15s
kubectl apply -f https://github.com/grafana/tempo-operator/releases/latest/download/tempo-operator.yaml
sleep 15s
kubectl apply -f https://raw.githubusercontent.com/grafana/tempo-operator/main/minio.yaml
kubectl apply -n monitoring -f .
