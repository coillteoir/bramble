#!/bin/sh

set -e

FAIL=0

kubectl apply -k tests/e2e/resources/

echo "=======SLEEPING FOR 60 SECONDS======="
sleep 60

echo "RUNNING TEST: SIMPLE PIPELINE"

SIMPLE_PODS=$(kubectl -n testns get po \
-l=bramble-execution=simple-test \
-o jsonpath='{.items[*].metadata.name}')

echo "$SIMPLE_PODS"

if [ $(echo "$SIMPLE_PODS" | wc -w) -eq '2' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi


echo "RUNNING TEST: SLEEPY PIPELINE"

SLEEPY_PODS=$(kubectl -n testns get po \
-l=bramble-execution=sleepytest \
-o jsonpath='{.items[*].metadata.name}')

echo "$SLEEPY_PODS"

if [ $(echo "$SLEEPY_PODS" | wc -w) -eq '7' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi

echo "RUNNING TEST: ONE-TO-MANY PIPELINE"

ONE_TO_MANY_PODS=$(kubectl -n testns get po \
-l=bramble-execution=one-to-many \
-o jsonpath='{.items[*].metadata.name}')

echo "$ONE_TO_MANY_PODS"

if [ $(echo "$ONE_TO_MANY_PODS" | wc -w) -eq '6' ] 
    then echo "Correct amount of pods have been created." 
    else echo "Incorrect amount of pods created." && FAIL=1
fi

kubectl delete -k tests/e2e/resources/
exit $FAIL
