#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

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
		delphix_engine_version = "${DELPHIX_VERSION}"
	EOF
	terraform init
	if [[ "${1}" == "redeploy" ]]; then
		shift
		terraform destroy "$@"
		terraform apply "$@"
	else
		terraform "$@"
	fi
  ENDTIME
}
