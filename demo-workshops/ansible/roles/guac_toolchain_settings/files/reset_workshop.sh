#!/bin/bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
# This file is managed by Ansible. Don't make changes here - they will be overwritten.


STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")

function ENDTIME {
	ENDTIME=$(date +%s)
	echo "It took $(($ENDTIME - $STARTTIME)) seconds to complete ${SCRIPT}"
}

function ERROR {
	ENDTIME
	mv ~/Desktop/WAIT ~/Desktop/ERROR
	exit 1
}

function READY {
	ENDTIME=$(date +%s)
	echo "Script finished Successfully"
	mv /home/ubuntu/Desktop/WAIT /home/ubuntu/Desktop/READY
}

rm -Rf ~/Desktop/READY ~/Desktop/WAIT ~/Desktop/ERROR

{
	rm -Rf ~/git/app_repo

	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	ssh centos@tooling <<-EOM
	cd /var/lib/jenkins
	sudo rm -Rf app_repo datical
	sudo tar xvf app_repo.tgz
	sudo tar xvf datical.tgz
	sudo chown -R centos.jenkins app_repo.git
	sudo chmod -R g+rwX app_repo.git
	sudo find app_repo.git -type d -exec chmod g+s '{}' +
	sudo chown -R centos.jenkins datical
	sudo chmod -R g+rwX datical
	sudo find datical -type d -exec chmod g+s '{}' +
	
	EOM
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	ssh -t centos@proddb "sudo su - oracle -c 'cd /home/oracle/patients && ./reset_patients.sh'"
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	ssh -t centos@tooling java -jar /opt/jenkins-cli.jar -auth admin:admin -s http://tooling:8080/ build -s -v PatientsPipeline/production
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	curl -v --retry 22 --retry-delay 5  --retry-connrefused http://prodweb:8080/auth/sign-up -H 'Content-Type: application/json' -H 'cache-control:
	no-cache' -d '{
			"username": "delphix_admin",
			"firstname": "Delphix",
			"lastname": "Admin",
			"password": "delphix"
	}' 
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	ssh -t centos@tooling snap_prod_refresh_mm --config snap_conf.txt 
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	ssh -t centos@tooling java -jar /opt/jenkins-cli.jar -auth admin:admin -s http://tooling:8080/ build -s -v PatientsPipeline/master &
	ssh -t centos@tooling java -jar /opt/jenkins-cli.jar -auth admin:admin -s http://tooling:8080/ build -s -v PatientsPipeline/develop &
	for job in `jobs -p`
	do
	echo $job
	wait $job || let "FAIL+=1"
	done

	[[ -n "$FAIL" ]] && echo "script failed" && exit 1

	rm -Rf ~/git/app_repo
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	git clone centos@tooling:/var/lib/jenkins/app_repo.git git/app_repo
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	cd git/app_repo
	git checkout develop
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	rsync -arv ~/notes_changes/* ~/git/app_repo/
	git add -A
	git commit -m "Addng Notes Field"
	[[ ${PIPESTATUS[0]} -ne 0 ]] && ERROR

	READY
} 2>&1 | tee ~/Desktop/WAIT
