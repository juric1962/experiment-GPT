# Minikube load test setup with Locust

## 1) Build and load server image in Minikube

Use the repository root as the build context so Docker can read `go.mod`, `go.sum`, and `main.go`.

```bash
eval $(minikube docker-env)
docker build -t experiment-gpt:latest -f Dockerfile .
```

## 2) Apply manifests

```bash
kubectl apply -f k8s/server-deployment-service.yaml
kubectl apply -f k8s/locust-configmap.yaml
kubectl apply -f k8s/locust-deployment-service.yaml
```

## 3) Open Locust UI

```bash
minikube service locust-web --url
```

Then open the URL and run the load test. Target host is already set to `http://experiment-gpt-service:9999`.
