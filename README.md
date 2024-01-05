# Bramble - A Kubernetes Native CI/CD Framework

## Installation
To install and run Bramble, you will need priveliged access to a Kubernetes cluster. 
For installation on your development machine I recommend Kind.

``` sh
kind create cluster
make local_k8s_deploy
# if you wish to run some demo pipelines then run the following
kubectl apply -k operator/config/samples
``` 

### Build dependencies:
- docker
- make

## Developer Experience
- Pipelines as Kubernetes manifests
- Pipelines can have handmade tasks or plug and play pre-applied tasks

### Demo ideas
- Pipelines to build, test, and deploy different components of Lugh
- C program to generate pngs, merging to master posts the generated image to Instagram