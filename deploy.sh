#!/usr/bin/env bash

set -exu

function deployDir() {
	cd "$1"
	./deploy.sh
	cd ..
}

kubectl apply -f .

deployDir ./grafana/
deployDir ./kube-state-metrics/
deployDir ./prometheus/
deployDir ./node-exporter/

# Deploy in order, wait for previous to complete before passing on to the next:

sleep 20s
deployDir ./minio/
deployDir ./alloy
sleep 20s
deployDir ./loki
sleep 20s
deployDir ./tempo/
sleep 20s
deployDir ./dice-app/
# sleep 20s
# deployDir ./demo-app/
