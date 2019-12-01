#!/usr/bin/env bash
set -e
echo "###FIX NOTES###"


if [[ $(ssh centos@proddb sudo -i -u postgres whoami) ]]; then
    cat > add_notes.sql <<-\EOF
ALTER TABLE PATIENTS
  ADD NOTES character varying(2000);
EOF
else
    cat > add_notes.sql <<-\EOF
ALTER TABLE PATIENTS
  ADD NOTES varchar2(2000);
EOF
fi

scp add_notes.sql ubuntu@jumpbox:git/app_repo/sql_code/ddl/.

rm add_notes.sql

ssh ubuntu@jumpbox 'cd git/app_repo && git add -A && git commit -m "corrected notes column" && git push'

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | jq -r .result| grep null > /dev/null || GREP_RETURN_CODE=1

done

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "SUCCESS" ]]; then
    echo "Patients Pipeline executed successfully"
else
    echo "Patients Pipeline Failed"
    exit 1
fi

echo "Testing the app for the presence of the notes field"

URL="http://devweb:8080"

AUTH=$(curl -sX POST \
  ${URL}/auth/login \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 2b89c369-7e87-4894-95f8-f3eeae99119f' \
  -H 'cache-control: no-cache' \
  -d '{
  "username": "patients_admin",
  "password": "delphix"
}')

if [[ $(curl ${URL}/patients/1   -H 'Content-Type: application/json'   -H "Authorization: Bearer ${AUTH}"   -H 'cache-control: no-cache'| jq -r 'has("notes")') ]]; then
  echo "Notes Field Successfully Added"
else
  echo "Notes Field not Present"
fi