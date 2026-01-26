#!/usr/bin/env bash

set -exu

kubectl apply -f .
kubectl apply -f ./grafana/
kubectl apply -f ./kube-state-metrics/
kubectl apply -f ./prometheus/
kubectl apply -f ./node-exporter/
