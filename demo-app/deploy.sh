#!/usr/bin/env bash

set -exu

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm install otel-demo open-telemetry/opentelemetry-demo -n demo-app --create-namespace --values ./values.yml
