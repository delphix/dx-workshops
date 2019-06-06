#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT
STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")
WORKDIR=$(pwd)

{
	GET_CLEANUP_LIST ${TEMPLATE_LIST[@]}
	[[ -z ${CLEANUP_LIST} ]] && echo "No AMI's to cleanup" && exit 0
	until [[ "${CLEANUP}" == "y" || "${CLEANUP}" == "n" ]]; do
      read -p "This will delete all the AMI's listed above. Are you sure you want to continue? (n) " CLEANUP
      CLEANUP=${CLEANUP:-n}
    done
    if [[ "${CLEANUP}" == "y" ]]; then
      CLEANUP_AMIS $CLEANUP_LIST
    else
      echo "Cleanup cancelled"
		fi
}
