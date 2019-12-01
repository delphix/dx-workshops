#!/usr/bin/env bash
set -e
echo "###MERGE NOTES###"

./merge_feature.sh

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/master/lastBuild/api/json

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "SUCCESS" ]]; then
   echo "Job passed"
else
    echo "Patients Pipeline Didn't Succeed"
    exit 1
fi