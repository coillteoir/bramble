#!/bin/sh

set -e

FAIL=0
SLEEPTIME=300

kubectl apply -k tests/e2e/resources/

echo "=======SLEEPING FOR $SLEEPTIME SECONDS======="
sleep $SLEEPTIME

for execution in $(kubectl -n testns get executions -o jsonpath="{.items[*].metadata.name}"); do
    echo RUNNING TEST: "$execution"
    jobcount="$(kubectl -n testns get jobs \
    -l=bramble-execution="$execution" \
    -o jsonpath='{.items}' | jq length)"

    taskcount="$(kubectl -n testns get pipeline \
        "$(kubectl -n testns get execution "$execution" -o jsonpath="{.spec.pipeline}")"\
        -o jsonpath="{.spec.tasks}" | jq length)"

    if [ "$(echo "$jobcount" | wc -w)" -eq "$taskcount" ]
        then echo Correct amount of jobs have been created in "$execution".
        else echo Incorrect amount of jobs created in "$execution". && FAIL=1
    fi
done

kubectl delete -k tests/e2e/resources/
exit $FAIL
