#!/usr/bin/env bash
set -e
echo "###EXECUTE DATA PIPELINE###"
ssh -o "StrictHostKeyChecking=no" centos@delphix-tcw-tooling-postgres11 java -jar /opt/jenkins-cli.jar -s http://localhost:8080 -auth admin:admin build \"Data Pipeline\"

JOB_STATUS_URL=http://delphix-tcw-tooling-postgres11:8080/job/Data%20Pipeline/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for Data Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | grep result\":null > /dev/null || GREP_RETURN_CODE=1
done

curl --silent $JOB_STATUS_URL