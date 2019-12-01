#!/usr/bin/env bash
set -e
echo "###FIX TEST###"

if [[ $(ssh centos@proddb sudo -i -u postgres whoami) ]]; then
    cat > add_constraint.sql <<-\EOF
ALTER TABLE ONLY public.users 
        ADD CONSTRAINT username UNIQUE(username);
EOF
else
    cat > add_constraint.sql <<-\EOF
ALTER TABLE USERS
    ADD CONSTRAINT username UNIQUE(username);
EOF
fi

scp add_constraint.sql ubuntu@jumpbox:git/app_repo/sql_code/ddl/.

rm add_constraint.sql

ssh ubuntu@jumpbox 'cd git/app_repo && git checkout develop && git pull && git add -A && git commit -m "added constraint" && git push'

JOB_STATUS_URL=http://tooling:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null |  jq -r .result| grep null > /dev/null || GREP_RETURN_CODE=1
done

if [[ "$(curl -s $JOB_STATUS_URL| jq -r .result)" == "SUCCESS" ]]; then
    echo "Patients Pipeline executed successfully"
else
    echo "Patients Pipeline Failed"
    exit 1
fi