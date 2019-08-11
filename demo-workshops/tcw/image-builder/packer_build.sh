#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "packer_cleanup" SIGINT

: "${VNC_CLIENT_OPTIONS:=-geometry 1280x720 -localhost yes}"
export VNC_CLIENT_OPTIONS

function packer_cleanup() {
	echo "Caught CTRL+C. Terminating packer jobs"
	for child in $(ps aux| grep '[/]bin/packer build' | awk '{print $1}' ); do
		echo kill "$child" && kill -s SIGINT "$child"
	done
	wait $(jobs -p)
	echo "You may need to go in and manually terminate instances and delete security groups and keypairs (search for items with 'packer' in the name)"
  	ERROR
}

rm -f READY.log WAIT.log ERROR.log change.ignore

function ALL() {
	for BATCH in `GET_BATCHES`; do
		echo "BATCH ${BATCH}"
		PACKER_BUILD ${BATCH}
	done
}

{
	ENVCHECK
	CERT_TEST
	BINARY_BUILD
	if [[ -z "$1" ]]; then
		ALL
	else
		echo "Starting Batch ${1}"
		PACKER_BUILD ${1}
	fi
} 2>&1 | tee ${WORKDIR}/WAIT.log
READY