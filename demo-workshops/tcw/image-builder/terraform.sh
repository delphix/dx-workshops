#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
trap "cleanup" SIGINT
STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")
WORKDIR=$(pwd)
DEMO_PATH="demo-workshops"
DEMO_NAME="tcw"
TERRAFORM_BLUEPRINTS="${WORKDIR}/${DEMO_PATH}/${DEMO_NAME}/terraform-blueprints"

function cleanup() {
	echo "Caught CTRL+C. Terminating packer jobs"
	for child in $(ps aux| grep '[/]bin/packer build' | awk '{print $1}' ); do
		echo kill "$child" && kill -s SIGINT "$child"
	done
	wait $(jobs -p)
	echo "You may need to go in and manually terminate instances and delete security groups and keypairs (search for items with 'packer' in the name)"
}

function ENDTIME {
	ENDTIME=$(date +%s)
	echo "It took $(($ENDTIME - $STARTTIME)) seconds to complete ${0}"
}

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
	EOF
	terraform init
	terraform "$@"
  ENDTIME
}
