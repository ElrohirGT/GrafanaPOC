#!/usr/bin/env bash

set -exu
# Load image into cluster
docker build -t dice-app .
CONTEXT=$(kubectl config current-context 2>/dev/null || true)
if [[ "$CONTEXT" == *"kind"* ]]; then
  echo "📥 Loading image into kind cluster..."
  kind load docker-image dice-app:latest
elif [[ "$CONTEXT" == *"k3d"* ]]; then
  echo "📥 Loading image into k3d cluster..."
  K3D_CLUSTER="${CONTEXT#k3d-}"
  k3d image import dice-app:latest -c "$K3D_CLUSTER"
elif command -v minikube &>/dev/null && [[ "$CONTEXT" == *"minikube"* ]]; then
  echo "📥 Loading image into minikube..."
  minikube image load dice-app:latest
fi

kubectl apply -f ./kubernetes/
