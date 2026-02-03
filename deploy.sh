#!/usr/bin/env bash

set -exu

kubectl apply -f .
kubectl apply -f ./grafana/
kubectl apply -f ./kube-state-metrics/
kubectl apply -f ./prometheus/
kubectl apply -f ./node-exporter/

function deployDir() {
	cd "$1"
	./deploy.sh
	cd ..
}
#
# # Deploy in order, wait for previous to complete before passing on to the next:
#
deployDir ./loki
sleep 20s
deployDir ./alloy
sleep 20s
deployDir ./tempo/
sleep 20s
deployDir ./demo-app/
# # Not required
# deployDir ./jaeger
