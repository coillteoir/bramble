#!/bin/sh

set -e

FAIL=0
SLEEPTIME=300

kubectl apply -k tests/e2e/resources/

echo "=======SLEEPING FOR $SLEEPTIME SECONDS======="
sleep $SLEEPTIME

for execution in $(kubectl -n testns get executions -o jsonpath="{.items[*].metadata.name}"); do
    echo RUNNING TEST: "$execution"
    podcount="$(kubectl -n testns get po \
    -l=bramble-execution="$execution" \
    -o jsonpath='{.items}' | jq length)"

    taskcount="$(kubectl -n testns get pipeline \
        "$(kubectl -n testns get execution "$execution" -o jsonpath="{.spec.pipeline}")"\
        -o jsonpath="{.spec.tasks}" | jq length)"

    if [ "$(echo "$podcount" | wc -w)" -eq "$taskcount" ]
        then echo Correct amount of pods have been created in "$execution".
        else echo Incorrect amount of pods created in "$execution". && FAIL=1
    fi
done

kubectl delete -k tests/e2e/resources/
exit $FAIL
