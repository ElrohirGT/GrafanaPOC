# Grafana

## Prerequisites:

- Kubernetes Cluster (Kind or another).
- Helm

## Running

First you need to start the pods, run the `./deploy.sh` script to initialize
prometheus and grafana.

To initialize `loki` and `alloy` run:

```bash
cd ./loki/ && helm install loki grafana/loki -n loki --create-namespace -f loki.yaml && cd ..
cd ./alloy/ && helm install grafana-alloy grafana/alloy -n loki -f alloy.yml && cd ..
```

To port forward and access Grafana type:

```bash
kubectl port-forward service/grafana 3000:3000 --namespace=test
```

## References

- https://devopscube.com/setup-grafana-loki/
- https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/
- https://devopscube.com/setup-grafana-kubernetes/
- https://devopscube.com/node-exporter-kubernetes/
- https://devopscube.com/setup-kube-state-metrics/
