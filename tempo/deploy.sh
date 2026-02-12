#!/usr/bin/env bash

set -exu

helm repo add grafana-community https://grafana-community.github.io/helm-charts
helm install tempo grafana-community/tempo-distributed -n monitoring --values ./values.yml
