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

{
  for each in "${SYSTEMS[@]}"; do
    SYSTEM_IDS+=($(terraform output ${each}))
  done

  case ${1} in
  start)
    aws ec2 --region ${AWS_REGION} start-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  stop)
    aws ec2 --region ${AWS_REGION} stop-instances --instance-ids ${SYSTEM_IDS[@]}
    ;;
  esac
}