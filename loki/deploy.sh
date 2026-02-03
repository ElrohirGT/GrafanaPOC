#!/usr/bin/env bash

set -exu

helm repo add grafana https://grafana.github.io/helm-charts
helm install loki grafana/loki -n monitoring --create-namespace -f loki.yaml
