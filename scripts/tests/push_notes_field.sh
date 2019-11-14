#!/usr/bin/env bash
set -e
echo "###PUSH NOTES FIELD###"

ssh -o "StrictHostKeyChecking=no" ubuntu@delphix-tcw-jumpbox 'cd git/app_repo && git commit -m "notes field" && git push'

JOB_STATUS_URL=http://delphix-tcw-tooling-postgres11:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | grep result\":null || GREP_RETURN_CODE=1
done

echo "Exit code"
curl --silent $JOB_STATUS_URL 2>/dev/null | grep add_notes.sql