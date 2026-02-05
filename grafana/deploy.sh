#!/usr/bin/env bash

set -exu

# Create configmap
# kubectl create configmap my-app-config \
#   --from-file=./dashboards/1860_rev42.json \
#   --from-file=./dashboards/13332_rev12.json \
#   --dry-run=client -o yaml > configmap.yaml
