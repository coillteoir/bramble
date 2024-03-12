#!/bin/sh

set -e

FAIL=0

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
    else echo "Incorrect amount of pods created." && FAIL=1
fi


echo "RUNNING TEST: SLEEPY PIPELINE"

SLEEPY_PODS=$(kubectl get po \
-l=bramble-execution=sleepytest \
-o jsonpath='{.items[*].metadata.name}')

echo "$SLEEPY_PODS"

if [ $(echo "$SLEEPY_PODS" | wc -w) -eq '7' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi

echo "RUNNING TEST: ONE-TO-MANY PIPELINE"

ONE_TO_MANY_PODS=$(kubectl get po \
-l=bramble-execution=one-to-many \
-o jsonpath='{.items[*].metadata.name}')

echo "$ONE_TO_MANY_PODS"

if [ $(echo "$ONE_TO_MANY_PODS" | wc -w) -eq '6' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi

echo "RUNNING TEST: ONE-TO-MANY PIPELINE"

ONE_TO_MANY_PODS=$(kubectl get po \
-l=bramble-execution=one-to-many \
-o jsonpath='{.items[*].metadata.name}')

echo "$ONE_TO_MANY_PODS"

if [ $(echo "$ONE_TO_MANY_PODS" | wc -w) -eq '6' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi

kubectl delete pipeline sleepy
kubectl delete execution sleepytest
kubectl delete pipeline simple-test
kubectl delete execution simple-test
kubectl delete pipeline one-to-many
kubectl delete execution one-to-many

exit $FAIL
