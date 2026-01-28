#!/usr/bin/env bash

set -exu

helm install grafana-alloy grafana/alloy -n loki -f alloy.yml
