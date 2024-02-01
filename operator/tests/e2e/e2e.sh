#!/bin/sh

set -e


kubectl apply -f tests/e2e/resources/one-to-many-execution.yaml
kubectl apply -f tests/e2e/resources/one-to-many-pipeline.yaml
kubectl apply -f tests/e2e/resources/simple_execution.yaml
kubectl apply -f tests/e2e/resources/simple_pipeline.yaml
kubectl apply -f tests/e2e/resources/sleepy_execution.yaml
kubectl apply -f tests/e2e/resources/sleepy_pipeline.yaml

echo "=======SLEEPING FOR 60 SECONDS======="
sleep 60

echo "RUNNING TEST: SIMPLE PIPELINE"

SIMPLE_PODS=$(kubectl get po \
-l=bramble-execution=simple-test \
-o jsonpath='{.items[*].metadata.name}')

echo "$SIMPLE_PODS"

if [ $(echo "$SIMPLE_PODS" | wc -w) -eq '2' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-test-cluster && false
fi


echo "RUNNING TEST: SLEEPY PIPELINE"

SLEEPY_PODS=$(kubectl get po \
-l=bramble-execution=sleepytest \
-o jsonpath='{.items[*].metadata.name}')

echo "$SLEEPY_PODS"

if [ $(echo "$SLEEPY_PODS" | wc -w) -eq '7' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-sleepy-cluster && false
fi

echo "RUNNING TEST: ONE-TO-MANY PIPELINE"

SLEEPY_PODS=$(kubectl get po \
-l=bramble-execution=one-to-many \
-o jsonpath='{.items[*].metadata.name}')

echo "$SLEEPY_PODS"

if [ $(echo "$SLEEPY_PODS" | wc -w) -eq '6' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && kind delete cluster --name=bramble-sleepy-cluster && false
fi

kubectl delete pipeline sleepy
kubectl delete execution sleepytest
kubectl delete pipeline simple
kubectl delete execution simple
kubectl delete pipeline one-to-many
kubectl delete execution one-to-many