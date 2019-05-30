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

SYSTEMS=(delphix-tcw-delphixengine_id delphix-tcw-jumpbox_id delphix-tcw-oracle12-source_id delphix-tcw-oracle12-target_id delphix-tcw-tooling-oracle_id devweb_id prodweb_id testweb_id)

SYSTEM_IDS=()
cd ${TERRAFORM_BLUEPRINTS}

function STATUS(){
  aws ec2 --region ${AWS_REGION} describe-instances --instance-ids ${SYSTEM_IDS[@]} --output text --query "Reservations[*].Instances[*].[ImageId,State.Name]"
}

function WAIT_FOR(){
  echo "Checking status of instances"
  if [[ "${1}" == "start" ]]; then
    STATE="running"
  else
    STATE="stopped"
  fi

  until [[ -z "`STATUS | grep -v ${STATE}`" ]]; do
    STATUS | grep -v ${STATE}
    echo "Will check again in 5 seconds"
    sleep 5
  done
  echo "All instances ${STATE}"
  STATUS
}

{
  for each in "${SYSTEMS[@]}"; do
    SYSTEM_IDS+=($(terraform output ${each}))
  done

  case ${1} in
  start)
    echo "Starting stopped instances"
    aws ec2 --region ${AWS_REGION} start-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  stop)
    echo "Stopping running instances"
    aws ec2 --region ${AWS_REGION} stop-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  esac

  [[ ${2} == "wait" ]] && WAIT_FOR ${1}
}