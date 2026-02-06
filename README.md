# Grafana

## Prerequisites:

- Kubernetes Cluster: k3d recommended but you can also use another.
- Helm

If you user k3d, to create the cluster simply run:

```bash
k3d cluster create
```

## Running

This cluster is heavy, so I recommend using [k3d](https://k3d.io/v5.6.3/). First
you need to start the pods, run the `./deploy.sh` script to initialize all
services. It's likely some deployments will fail due to missing dependencies,
please restart the deployments from within that folder once it does! Normally
just one restart is required!

It's recommended that you read the `./deploy.sh` script!

To port forward and access Grafana type:

```bash
kubectl port-forward service/grafana 3000:3000 --namespace=test
```

Then simply access the GrafanaUI from within a web browser:

```
http://localhost:3000
```

To port forward the dice-app type:

```bash
kubectl port-forward service/dice-service 8081:8080 --namespace=dice-app
```

Then you can make get requests like so:

```bash
curl http://localhost:8081/rolldice/Jose
```

## More complex example

If you need a more complex example, you can deploy the preconfigured `demo-app`.
This is the opentelemetry demo adapted to use the LGTM stack on this repo.

Simply run:

```bash
# From the repo root
cd ./demo-app/ && ./deploy.sh
```

## References

- https://devopscube.com/setup-grafana-loki/
- https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/
- https://devopscube.com/setup-grafana-kubernetes/
- https://devopscube.com/node-exporter-kubernetes/
- https://devopscube.com/setup-kube-state-metrics/
- https://grafana.com/docs/tempo/latest/set-up-for-tracing/setup-tempo/deploy/kubernetes/operator/quickstart/
- https://github.com/jaegertracing/jaeger-operator?tab=readme-ov-file#jager-v2-operator
