# Learn kubernetes project

Main aim of project - practical use of most Kubernetes features\functions\objects and tools.

Including, not but not limited to:

- Deploy
- Pod
- Service
- Secret
- EnvVariable
- PVC
- ingress
- RBAC
- jobs

Tools: kubectl, helm.

#### Requirements

Before install central services - install and config ingres.
To install nginx-ingres use command similar to:

```sh
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx
```

Get details at [https://kubernetes.github.io/ingress-nginx/deploy/](https://kubernetes.github.io/ingress-nginx/deploy/)

App using next ports (shouldn't be blocked):
8001 - service echo
8002 - service timer
8003 - service colorer

#### Install

```sh
helm install echo-cluster https://github.com/ups91/learn-kubernetes/raw/main/echo-chart-0.1.0.tgz
```
