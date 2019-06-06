#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
export DELPHIX_VERSION="5.3.2.*"
TEMPLATE_LIST=(delphix-centos7-ansible-base.json delphix-ubuntu-bionic-guacamole.json \
		delphix-centos7-oracle-12.2.0.1.json delphix-centos7-daf-app.json \
		delphix-centos7-kitchen_sink.json delphix-tcw-jumpbox.json delphix-tcw-oracle12-source.json \
		delphix-tcw-oracle12-target.json delphix-centos7-tooling-base.json delphix-tcw-tooling-oracle.json)
SYSTEMS=(delphix-tcw-delphixengine_id delphix-tcw-jumpbox_id delphix-tcw-oracle12-source_id delphix-tcw-oracle12-target_id delphix-tcw-tooling-oracle_id devweb_id prodweb_id testweb_id)

trap "cleanup" SIGINT

function cleanup() {
	echo "Caught CTRL+C. Terminating packer jobs"
	for child in $(ps aux| grep '[/]bin/packer build' | awk '{print $1}' ); do
		echo kill "$child" && kill -s SIGINT "$child"
	done
	wait $(jobs -p)
	echo "You may need to go in and manually terminate instances and delete security groups and keypairs (search for items with 'packer' in the name)"
}

function ENVCHECK() {
	[[ -z "${S3_AWS_ACCESS_KEY_ID}" ]] && export S3_AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
	[[ -z "${S3_AWS_SECRET_ACCESS_KEY}" ]] && export S3_AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
}

function ENDTIME() {
	ENDTIME=$(date +%s)
	echo "It took $(($ENDTIME - $STARTTIME)) seconds to complete ${0}"
}

function ERROR() {
	ENDTIME
	mv ${WORKDIR}/WAIT.log ${WORKDIR}/ERROR.log
	exit 1
}

function READY() {
	ENDTIME
	echo "Script finished Successfully"
	mv ${WORKDIR}/WAIT.log ${WORKDIR}/READY.log
}

function JOB_WAIT() {
	for job in `jobs -p`
	do
		wait $job || let "FAIL+=1"
	done
	[[ -n "${FAIL}" ]] && ERROR
  true #exits function cleanly, if no error condition
}

function CERT_TEST() {
	#Check to see if the automation certificates exist, if not prompt to create
	if [[ -f ${CERT}.pub && -f ${CERT} ]] ; then
		echo "Certifcate: ${CERT}"
	elif [[ ! -f ${CERT}.pub && -f ${CERT} ]]; then
    until [[ "${CREATEPUB}" == "y" || "${CREATEPUB}" == "n" ]]; do
      read -p "${CERT}.pub is missing. Should I create it? (y) " CREATEPUB
      CREATEPUB=${CREATEPUB:-y}
    done
    if [[ "${CREATEPUB}" == "y" ]]; then
      ssh-keygen -y -f ${CERT} > ${CERT}.pub
    else
      echo "set the CERT variable to the private certificate you would like to use for automation"
      echo "The public certificate should be found in the same directory"
      echo "This should be a key exclusive to run in this environment"
      exit 1
    fi
  else
    until [[ "${CREATECERT}" == "y" || "${CREATECERT}" == "n" ]]; do
      read -p "${CERT} is missing. Should I create it? (y) " CREATECERT
      CREATECERT=${CREATECERT:-y}
    done
    if [[ "${CREATECERT}" == "y" ]]; then
      ssh-keygen -b 2048 -t rsa -f certs/ansible -q -N ""
    else
      echo "set the CERT variable to the private certificate you would like to use for automation"
      echo "The public certificate should be found in the same directory"
      echo "This should be a key exclusive to run in this environment"
      exit 1
    fi
	fi
	export ANSIBLE_PUB_KEY=$(cat ${CERT}.pub)
}

function AMI_INFO() {
    for each in "$@"; do
        # query for an existing AMI with the name and md5sum number, and store that information in a file
        cd  $(RETURN_DIRECTORY $each)
        aws ec2 --region ${AWS_REGION} describe-images --filters "Name=name,Values=${each%.json}-*" "Name=tag:md5sum,Values=$(jq -r '.md5sum' ${each%.json}_md5sum.json)" --query 'sort_by(Images, &CreationDate)[-1]' > /tmp/${each%.json}_info.json
    done
    JOB_WAIT
}

function AMI_EXISTS() {
	# if a an ImageId is present, then the ami exists
    AMI=$(jq -r '.ImageId' /tmp/${1%.json}_info.json)
	[[ -n $AMI && $AMI != 'null' ]] 
}

function AMI_OLDER_THAN_SOURCE() {
	# Check the source tag of the AMI and compare the CreationDate of the AMI and it's Source, if applicable
	SOURCE_AMI="$(cat /tmp/${1%.json}_info.json | jq -r '.Tags[]|select(.Key=="source").Value')" 
    if [[ -f /tmp/${SOURCE_AMI%.json}_info.json ]]; then
        echo comparing ${1%.json}: $(cat /tmp/${1%.json}_info.json | jq -r '.CreationDate') to ${SOURCE_AMI%.json}: $(cat /tmp/${SOURCE_AMI%.json}_info.json | jq -r '.CreationDate')
        [[ "$(cat /tmp/${1%.json}_info.json | jq -r '.CreationDate')" < "$(cat /tmp/${SOURCE_AMI%.json}_info.json | jq -r '.CreationDate')" ]]
    else
        false
    fi
}

function NEED_TO_BUILD_AMI() {
  cd  $(RETURN_DIRECTORY $1)
  if AMI_EXISTS $1; then
      # If an AMI with a valid name and md5sum exists
      if AMI_OLDER_THAN_SOURCE $1; then
          # Check to see if it is older than it's source AMI. If so, build it
          echo ${1%.json} exists, but is older than source. Rebuilding.
          true
      else
          echo Skipping ${1}: a valid ami exists
          false
      fi
  else
      # Otherwise, build it
      echo Building ${1%.json}
      true
  fi
}

function PACKER_BUILD() {
    AMI_INFO "$@"
    for each in "$@"; do
      if NEED_TO_BUILD_AMI $each; then
        #drop this little file to alert our build server know a new ami was built
        touch ${WORKDIR}/change.ignore
        packer build -var-file ${each%.json}_md5sum.json $each &
      fi
    done
    JOB_WAIT
    AMI_INFO "$@"
}

function RETURN_DIRECTORY() {
	#This function is a placeholder until I figure how I want to reorg folder structure
	dirname $(find $WORKDIR -name $1)
}

function GET_CLEANUP_LIST() {
  for each in "$@"; do
    # query for an existing AMI with the name and md5sum number, and store that information in a file
    cd  $(RETURN_DIRECTORY $each)
    echo "Will deregister the following AMI's for ${each}:"
    for ami in $(aws ec2 --region ${AWS_REGION} describe-images --filters "Name=name,Values=${each%.json}-*" --query "Images[?Tags[?Key=='md5sum']|[?Value!='$(jq -r '.md5sum' ${each%.json}_md5sum.json)']].ImageId" --output text); do
      echo "${ami}"
      CLEANUP_LIST+="${ami} "
    done
    GET_OLDER_DUPLICATE_AMIS $each
  done
}

function CLEANUP_AMIS() {
  for each in "$@"; do
    echo "Deregistering ${each}"
    aws ec2 --region ${AWS_REGION} deregister-image --image-id ${each}
  done
}

function GET_OLDER_DUPLICATE_AMIS() {
  for each in "$@"; do
    cd  $(RETURN_DIRECTORY $each)
    for each in $(aws ec2 --region ${AWS_REGION} describe-images --filters "Name=name,Values=${each%.json}-*" "Name=tag:md5sum,Values=$(jq -r '.md5sum' ${each%.json}_md5sum.json)" --query "sort_by(Images, &CreationDate)[0:-1].ImageId" --output text); do
      echo "${each}"
      CLEANUP_LIST+="${each} "
    done
  done
}