#!/usr/bin/env bash
set -e
echo "###FIX TEST###"

cat > add_constraint.sql <<-\EOF
ALTER TABLE ONLY public.users 
    ADD CONSTRAINT username UNIQUE(username);
EOF

scp -o "StrictHostKeyChecking=no" add_constraint.sql ubuntu@delphix-tcw-jumpbox:git/app_repo/sql_code/ddl/.

rm add_constraint.sql

ssh -o "StrictHostKeyChecking=no" ubuntu@delphix-tcw-jumpbox 'cd git/app_repo && git checkout develop && git pull && git add -A && git commit -m "added constraint" && git push'

JOB_STATUS_URL=http://delphix-tcw-tooling-postgres11:8080/job/PatientsPipeline/job/develop/lastBuild/api/json

GREP_RETURN_CODE=0

while [ $GREP_RETURN_CODE -eq 0 ]
do
    sleep 5
    echo "waiting for develop Patients Pipeline" 
    
    curl --silent $JOB_STATUS_URL 2>/dev/null | grep result\":null > /dev/null || GREP_RETURN_CODE=1
done

curl --silent $JOB_STATUS_URL 2>/dev/null