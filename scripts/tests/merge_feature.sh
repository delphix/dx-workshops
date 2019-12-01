#!/usr/bin/env bash
set -e
echo "###MERGE FEATURE###"
ssh ubuntu@jumpbox 'cd git/app_repo && git pull && git checkout master && git merge develop && git push'

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/master/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for master Patients Pipeline" 
    
     curl --silent $JOB_STATUS_URL 2>/dev/null | jq -r .result| grep null > /dev/null || GREP_RETURN_CODE=1

done
