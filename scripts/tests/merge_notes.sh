#!/usr/bin/env bash
set -e
echo "###MERGE NOTES###"

./merge_feature.sh

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/master/lastBuild/api/json

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "UNSTABLE" ]]; then
    echo "Patients Pipeline expectedly failed"
    BASE_URL=$(curl -s $JOB_STATUS_URL | jq -r .url)
    TEST_STATUS=$(curl -s ${BASE_URL}testReport/junit/api/json?pretty=true | \
            jq -r ".suites[].cases[] | select(.className==\"Add Duplicate User\").status")

    if [[ "${TEST_STATUS}" != "REGRESSION" && "${TEST_STATUS}" != "FAILED" ]]; then
        echo "Did not have correct test failures"
        exit 1
    fi
    echo "Found correct results"
else
    echo "Patients Didn't go Unstable"
    exit 1
fi