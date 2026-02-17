#!/usr/bin/env bash

set -exu

helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install pyroscope grafana/pyroscope -n monitoring --wait
