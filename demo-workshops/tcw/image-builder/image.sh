#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

{ 
  [[ -z "$1" ]] && echo "Must provide a stage name to append to AMI. i.e. staged" && exit 1
  STAGE=${1}
  SHUTDOWN_VDBS
  cd ${TERRAFORM_BLUEPRINTS}
  for each in ${SYSTEMS[@]}; do
    if [[ $each == delphix-* ]]; then
      instance_id=$(terraform output -json ${each} | jq -r '.[]')
      AMI_NAME=${each%_id}-${STAGE}
      echo "${AMI_NAME}:${instance_id}"
      ansible-playbook -i 'localhost,' ${WORKDIR}/${DEMO_PATH}/ansible/ami_maker.yml -e "instance_id=${instance_id}" -e "ami_name=${AMI_NAME}" -e "commit=${GIT_COMMIT}" &
    fi
  done
  JOB_WAIT
  ENDTIME
}
