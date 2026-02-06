#!/usr/bin/env bash

set -exu

helm repo add grafana https://grafana.github.io/helm-charts
helm install grafana-alloy grafana/alloy -n monitoring -f alloy.yml || helm upgrade grafana-alloy grafana/alloy -n monitoring -f alloy.yml
