#!/usr/bin/env bash
set -e
echo "###EXECUTE DATA PIPELINE###"
ssh centos@tooling java -jar /opt/jenkins-cli.jar -s http://localhost:8080 -auth admin:admin build \"Data Pipeline\"

JOB_STATUS_URL=http://tooling:8080/job/Data%20Pipeline/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for Data Pipeline" 
    
    curl -s $JOB_STATUS_URL 2>/dev/null | jq -r .result| grep null > /dev/null || GREP_RETURN_CODE=1
done

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "SUCCESS" ]]; then
    echo "Patients Pipeline executed successfully"
else
    echo "Patients Pipeline Failed"
    exit 1
fi
