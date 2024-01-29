#!/bin/sh


which kind
which make
which go
which kubectl
which docker

set -e

kind create cluster --config=tests/e2e/test-cluster.yml

make deploy

echo "RUNNING TEST: SIMPLE PIPELINE"

kubectl apply -f tests/e2e/resources/simple_pipeline.yaml
kubectl apply -f tests/e2e/resources/simple_execution.yaml

echo "=======SLEEPING FOR 60 SECONDS======="
sleep 60


SIMPLE_PODS=$(kubectl get po \
-l=bramble-execution=simple-test \
-o jsonpath='{.items[*].metadata.name}')

echo "$SIMPLE_PODS"

if [ $(echo "$SIMPLE_PODS" | wc -w) -eq '2' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-test-cluster && false
fi

kubectl delete pipeline simple
kubectl delete execution simple

echo "RUNNING TEST: SLEEPY PIPELINE"

kubectl apply -f tests/e2e/resources/sleepy_pipeline.yaml
kubectl apply -f tests/e2e/resources/sleepy_execution.yaml

echo "=======SLEEPING FOR 60 SECONDS======="
sleep 60


SLEEPY_PODS=$(kubectl get po \
-l=bramble-execution=sleepytest \
-o jsonpath='{.items[*].metadata.name}')

echo "$SLEEPY_PODS"

if [ $(echo "$SLEEPY_PODS" | wc -w) -eq '7' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-sleepy-cluster && false
fi

kubectl delete pipeline sleepy
kubectl delete execution sleepytest

kind delete cluster --name=bramble-test-cluster
