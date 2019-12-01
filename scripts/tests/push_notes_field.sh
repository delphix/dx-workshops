#!/usr/bin/env bash
set -e
echo "###PUSH NOTES FIELD###"

ssh ubuntu@jumpbox 'cd git/app_repo && git commit -m "notes field" && git push'

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl -s $JOB_STATUS_URL 2>/dev/null | jq -r .result| grep null > /dev/null || GREP_RETURN_CODE=1
done

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "FAILURE" ]]; then
    echo "Patients Pipeline expectedly failed"
    PACKAGER_LOG=$(curl -s $JOB_STATUS_URL 2>/dev/null | jq -r ".artifacts[] | select(.displayPath==\"packager.log\").relativePath")
    BASE_URL=$(curl -s $JOB_STATUS_URL 2>/dev/null | jq -r ".url")
    if [[ ! $(curl -s ${BASE_URL}/artifact/${PACKAGER_LOG} | \
        grep "add_notes.sql failed validation") ]]; then
        echo "Did not have correct values in ${BASE_URL}/artifact/${PACKAGER_LOG}"
        exit 1
    fi
    echo "Found correct results"
else
    echo "Patients Didn't fail"
    exit 1
fi
