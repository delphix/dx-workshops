#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

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
  [[ "${1}" == "start" ]] && terraform refresh
}

{
  for each in $(terraform output  | grep _system | awk '{print $1}') ; do
    instance_id=$(terraform output -json ${each} | jq -r '.[]')
    echo "${each}:${instance_id}"
    SYSTEM_IDS+=( ${instance_id} )
  done

  case ${1} in
  start)
    echo "Starting stopped instances"
    aws ec2 --region ${AWS_REGION} start-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  stop)
    echo "Stopping running instances"
    SHUTDOWN_VDBS
    aws ec2 --region ${AWS_REGION} stop-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  esac

  [[ ${2} == "wait" ]] && WAIT_FOR ${1}
}