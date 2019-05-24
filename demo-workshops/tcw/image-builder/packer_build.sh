#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")
WORKDIR=$(pwd)
BASE_TEMPLATES="${WORKDIR}/base-templates"
DEMO_PATH="demo-workshops"
DEMO_NAME="tcw"
DEMO_TEMPLATES="${WORKDIR}/${DEMO_PATH}/${DEMO_NAME}/packer-templates"
CERT="certs/ansible"
GUACAMOLE_VERSION="0.9.14"
VNC_CLIENT_OPTIONS="-geometry 1280x720 -localhost"

rm -f READY.log WAIT.log ERROR.log change.ignore

function BATCH1() {
	echo "Starting Batch 1"
	PACKER_BUILD delphix-centos7-ansible-base.json delphix-ubuntu-bionic-guacamole.json
}

function BATCH2() {
	echo "Starting Batch 2"
	PACKER_BUILD delphix-centos7-oracle-12.2.0.1.json delphix-centos7-daf-app.json delphix-centos7-kitchen_sink.json delphix-tcw-jumpbox.json
}

function BATCH3() {
	echo "Starting Batch 3"
	PACKER_BUILD delphix-tcw-oracle12-source.json delphix-tcw-oracle12-target.json delphix-centos7-tooling-base.json
}

function BATCH4() {
	echo "Starting Batch 4"
	PACKER_BUILD delphix-tcw-tooling-oracle.json
}

function ALL() {
	BATCH1
	BATCH2
	BATCH3
	BATCH4
}

{
	ENVCHECK
  CERT_TEST
	if [[ -z "$1" ]]; then
		ALL
	else
		case ${1} in
			BATCH2)
				BATCH2
				BATCH3
				BATCH4
				;;
			BATCH3)
				BATCH3
				BATCH4
				;;
			BATCH4)
				BATCH4
				;;
			*)
				echo $1 is not a valid choice
				exit 1
				;;
			esac
	fi
} 2>&1 | tee ${WORKDIR}/WAIT.log
READY