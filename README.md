# Bramble - A Full Stack Kubernetes Native CI/CD System

Bramble is a CI/CD system where pipelines are designed to be run in parallel, with explicit dependency trees.
It utilizes the features of Kubernetes to run pipelines in parallel, tasks in isolation, and utilizing the scheduling provided by Kubernetes.

## Installation

To install and run Bramble, you will need priveliged access to a Kubernetes cluster.
For installation on your development machine I recommend Kind.

```sh
kind create cluster
make local_k8s_deploy
# if you wish to run some demo pipelines then run the following
kubectl apply -k operator/config/samples
```

### Build dependencies:

- kustomize
- kubectl
- docker
- make

### Dev dependencies

- docker
- make
- go
- gofumpt
- golangci-lint
- nodejs
- kustomize

## Developer Experience

- Pipelines as Kubernetes manifests.
- Pipelines can have handmade tasks or plug and play pre-applied tasks.
- Comprehensive UI for creating, executing and monitoring builds.
- Webhook support for Git providers.

### Demo ideas

- Pipelines to build, test, and deploy different components of Bramble
- C program to generate pngs, merging to master posts the generated image to instagram
- CI/CD pipeline for personal projects like MorningBot.
