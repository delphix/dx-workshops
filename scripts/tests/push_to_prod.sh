#!/usr/bin/env bash
set -e
echo "###PUSH TO PROD###"

ssh -o "StrictHostKeyChecking=no" ubuntu@guacamole 'cd git/app_repo && git pull && git checkout production && git merge master && git push'

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/production/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for production Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | grep result\":null > /dev/null || GREP_RETURN_CODE=1
done

curl --silent $JOB_STATUS_URL 2>/dev/null 