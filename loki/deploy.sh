#!/usr/bin/env bash

set -exu

helm install loki grafana/loki -n loki --create-namespace -f loki.yaml
