# Grafana

Once you have a local kubernetes instance running (EG: using kind). Run the
following command to deploy:

```bash
kubectl apply -f ./grafana-deployment.yml
```

If you do not see any EXTERNAL_IP, checkout:

```bash
kubectl port-forward service/grafana 3000:3000 --namespace=test
```
