#!/bin/sh


which kind
which make
which go
which kubectl
which docker

set -e

kind create cluster --config=tests/e2e/test-cluster.yml

make deploy

kubectl apply -f config/samples/sleepy_pipeline.yaml
kubectl apply -f config/samples/sleepy_execution.yaml

echo "=======SLEEPING FOR 60 SECONDS======="

sleep 60

PODNAMES=$(kubectl get po \
-l=bramble-execution=sleepytest \
-o jsonpath='{.items[*].metadata.name}')

if [ $(echo "$PODNAMES" | wc -w) -eq '7' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-sleepy-cluster && false
fi

kubectl delete pipeline sleepy
kubectl delete execution sleepytest

kind delete cluster --name=bramble-test-cluster
