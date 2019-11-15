#!/usr/bin/env bash
set -e
echo "###FIX NOTES###"

cat > add_notes.sql <<-\EOF
ALTER TABLE PATIENTS
  ADD NOTES character varying(2000);
EOF

scp -o "StrictHostKeyChecking=no" add_notes.sql ubuntu@guacamole:git/app_repo/sql_code/ddl/.

rm add_notes.sql

ssh -o "StrictHostKeyChecking=no" ubuntu@guacamole 'cd git/app_repo && git add -A && git commit -m "corrected notes column" && git push'

JOB_STATUS_URL=http://dtooling:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | grep result\":null > /dev/null || GREP_RETURN_CODE=1

done

curl --silent $JOB_STATUS_URL 2>/dev/null