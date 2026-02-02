#!/usr/bin/env bash

set -exu

helm install loki grafana/loki -n monitoring --create-namespace -f loki.yaml
