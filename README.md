# Grafana

## Prerequisites:

- Kubernetes Cluster (Kind or another).
- Helm

## Running

First you need to start the pods, run the `./deploy.sh` script to initialize
prometheus and grafana.

Read the script since some deployments need to be done manually.

To port forward and access Grafana type:

```bash
kubectl port-forward service/grafana 3000:3000 --namespace=test
```

Then simply access the GrafanaUI from within a web browser:

```
http://localhost:3000
```

To access the JaggerUI to see traces type:

```bash
kubectl port-forward svc/tempo-simplest-query-frontend 16686:16686
```

## References

- https://devopscube.com/setup-grafana-loki/
- https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/
- https://devopscube.com/setup-grafana-kubernetes/
- https://devopscube.com/node-exporter-kubernetes/
- https://devopscube.com/setup-kube-state-metrics/
- https://grafana.com/docs/tempo/latest/set-up-for-tracing/setup-tempo/deploy/kubernetes/operator/quickstart/
- https://github.com/jaegertracing/jaeger-operator?tab=readme-ov-file#jager-v2-operator
