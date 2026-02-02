#!/usr/bin/env bash

set -exu

helm install grafana-alloy grafana/alloy -n monitoring -f alloy.yml || helm upgrade grafana-alloy grafana/alloy -n monitoring -f alloy.yml
