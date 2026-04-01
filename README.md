# Ride-Sharing Microservices Platform

This repository contains a ride-sharing platform built with a microservices architecture. The current workspace includes Go services, a Next.js web client, local Kubernetes development with Tilt, and production deployment assets for selected services.

## Architecture

### Trip Scheduling Flow
[![](https://mermaid.ink/img/pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ?type=png)](https://mermaid.live/edit#pako:eNqNVt9v2jAQ_lcsP21qGvGjZSEPlSpaTX1YxWDVpAmpMvZBIkicOQ6UVf3fd4mdEgcKzQOK47v7vjt_d-aVcimAhnSW5vC3gJTDXcyWiiWzlOCTMaVjHmcs1eQpB3X49Xb88J1p2LLd4d4vFWdTUJuYw-HmnYo3oM5sH34fs10Cqf7Qb6oRFL-bnZLz5c3NnmRIRgrwteJGJmXOuTa2eyP0aFAPyXIyHtWO5YaxN7-PEoNJpNrM1nOSCw3Y_QuPWLq0nBvWl4h30fIos_Bhg5n6vMIVDTgVLyNN5II4LO9L65A8wrbyJimAyAkjwlbSBHBwENisQ2vl80T4pfezapbGRW1RHckkYaloILMNi9dsvgYXtHUQv2E-lXwF-hCccQ5ZE3sNiwZ0A_O2sjSw6oPTbB_nWbT3TB3dGEiykGrLlABBtCQTNp_H-sdP8sUwK3EmkGcS2-lnAQV8rUtwVCchGSvJIc8tJ2KoMGxDjxSZKIVapZZrpovc9_2j4nF7whGPifvM8jxepp8XkW2SzAQmOVKMZXoyF6_Nwq5bunetkLxp2HfIUQR8JQts5Bqz9DJGx3K1ZtZdHAVpCc9mZStkc3s-0WZtTFukOsGnB5QeE7u6PM4gKScQsozktmF_TmuN1qidUHZJa6q9V04m2RowlrV1SnbQc5GUq33YacFL_R2faHtPzxXt0ZN1O-6iXbRW1Zu4HxeiqjTJivk6ziPTcp-U1aXD-NPgZ46a21KL061Qj6kzc38_facarzCyv1tOdGgVk7OUpCipOSxjZEI9moBKWCzwKn8tQ8yojiCBGQ3xVTC1muEV_4Z2rNByuks5DbUqwKNKFsuIhgu2znFlZo79C1Cb4L36R8rmkoav9IWGvW_-1XVn0O_1-kE3GAyHgUd3-Lnb8fu9frc_xKfbvQ6CN4_-qyJ0_KDX7Q86QTDoDAfD66ve23_1IPGQ)

## Repository Components

- `services/apigateway`: HTTP and WebSocket gateway
- `services/trip-service`: trip orchestration and trip domain logic
- `services/driver-service`: driver-related backend service
- `web`: Next.js frontend
- `infra/deploy/development`: local Kubernetes and Docker assets used by Tilt
- `infra/deploy/production`: production Dockerfiles and Kubernetes manifests currently available in this repository

## Prerequisites

Install the following tools before starting the local environment:

- Docker
- Go
- Tilt
- `kubectl`
- A local Kubernetes cluster such as Minikube or Docker Desktop Kubernetes

### macOS

1. Install [Homebrew](https://brew.sh/).
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/).
3. Install [Minikube](https://minikube.sigs.k8s.io/docs/).
4. Install [Tilt](https://tilt.dev/).
5. Install Go:

```bash
brew install go
```

6. Install [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/).

### Windows (WSL)

Use WSL for the local development workflow.

1. Install [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).
2. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/).
3. Install [Minikube](https://minikube.sigs.k8s.io/docs/).
4. Install [Tilt](https://tilt.dev/).
5. Install Go inside WSL:

```bash
wget https://dl.google.com/go/go1.23.0.linux-amd64.tar.gz
sudo tar -xvf go1.23.0.linux-amd64.tar.gz
sudo mv go /usr/local
```

6. Add Go to your shell configuration:

```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

7. Reload your shell and verify the installation:

```bash
go version
```

8. Install [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/).

## Local Development

Start the development environment with Tilt:

```bash
tilt up
```

Monitor workloads with:

```bash
kubectl get pods
```

Or open the Minikube dashboard:

```bash
minikube dashboard
```

## Production Deployment Example (Google Cloud)

This section describes a manual deployment flow to Google Cloud. It is best used as a reference before formalizing the process in CI/CD.

### 1. Set environment variables

```bash
export REGION="europe-west1"
export PROJECT_ID="<your-gcp-project-id>"
```

### 2. Build production images

The production Dockerfiles currently present in this repository are:

```bash
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/ride-sharing/api-gateway:latest --platform linux/amd64 -f infra/deploy/production/docker/api-gateway.Dockerfile .
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/ride-sharing/trip-service:latest --platform linux/amd64 -f infra/deploy/production/docker/trip-service.Dockerfile .
```

### 3. Create an Artifact Registry repository

Create a Docker repository in Google Cloud Artifact Registry to store the built images.

### 4. Push the images

```bash
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/ride-sharing/api-gateway:latest
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/ride-sharing/trip-service:latest
```

If authentication fails:

1. Run `gcloud auth login` and select the correct project.
2. Configure Docker authentication:

```bash
gcloud auth configure-docker ${REGION}-docker.pkg.dev
```

### 5. Create a GKE cluster

Create a Google Kubernetes Engine cluster either through the `gcloud` CLI or the Google Cloud Console.

### 6. Connect to the cluster

```bash
gcloud container clusters get-credentials ride-sharing --region ${REGION} --project ${PROJECT_ID}
```

### 7. Update image references in the manifests

The production manifests under `infra/deploy/production/k8s` use `{{PROJECT_ID}}` placeholders. Replace those placeholders before applying the manifests.

### 8. Apply the manifests

```bash
kubectl apply -f infra/deploy/production/k8s/app-config.yaml
kubectl apply -f infra/deploy/production/k8s/trip-service-deployment.yaml
kubectl apply -f infra/deploy/production/k8s/api-gateway-deployment.yaml
```

For redeployments:

```bash
kubectl apply -f infra/deploy/production/k8s
kubectl rollout restart deployment
```

### 9. Verify service exposure

Retrieve the external address for the API gateway:

```bash
kubectl get services
```

To switch back to your local Kubernetes context:

```bash
kubectl config get-contexts

# Docker Desktop
kubectl config use-context docker-desktop

# Minikube
kubectl config use-context minikube
```
