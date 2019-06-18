#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT
STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")
WORKDIR=$(pwd)
DEMO_PATH="demo-workshops"
DEMO_NAME="tcw"
TERRAFORM_BLUEPRINTS="${WORKDIR}/${DEMO_PATH}/${DEMO_NAME}/terraform-blueprints"
{ 
  [[ -z "$1" ]] && echo "Must provide a stage name to append to AMI. i.e. staged" && exit 1
  STAGE=${1}
  SHUTDOWN_VDBS
  cd ${TERRAFORM_BLUEPRINTS}
  for each in ${SYSTEMS[@]}; do
    if [[ $each == delphix-* ]]; then
      instance_id=$(terraform output ${each})
      AMI_NAME=${each%_id}-${STAGE}
      echo "${AMI_NAME}:${instance_id}"
      ansible-playbook -i 'localhost,' ${WORKDIR}/${DEMO_PATH}/ansible/ami_maker.yml -e "instance_id=${instance_id}" -e "ami_name=${AMI_NAME}" -e "commit=${GIT_COMMIT}" &
    fi
  done
  JOB_WAIT
  ENDTIME
}
