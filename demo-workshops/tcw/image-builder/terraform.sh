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
  cd ${TERRAFORM_BLUEPRINTS}
	[[ -z ${AWS_AZ} ]] && AWS_AZ="${AWS_REGION}a"
	cat > terraform.tfvars <<- EOF
		access_key="${AWS_ACCESS_KEY_ID}"
		secret_key="${AWS_SECRET_ACCESS_KEY}"
		key_name="${AWS_KEYNAME}"
		project = "${AWS_PROJECT}"
		how_many = 1
		aws_region= "${AWS_REGION}"
		availability_zone = "${AWS_AZ}"
		delphix_engine_version = "Delphix Engine ${DELPHIX_VERSION}"
	EOF
	terraform init
	terraform "$@"
  ENDTIME
}
