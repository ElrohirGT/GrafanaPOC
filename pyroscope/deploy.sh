#!/usr/bin/env bash

set -exu

helm repo add grafana https://grafana.github.io/helm-charts
helm -n monitoring install pyroscope grafana/pyroscope
