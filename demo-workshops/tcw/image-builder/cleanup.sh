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
  if [[ -z "${1}" || "${1,,}" == "unstaged" ]]; then
	  GET_CLEANUP_LIST ${TEMPLATE_LIST[@]}
  elif [[ "${1,,}" == "all" ]]; then
    echo "WARNING: THIS WILL REMOVE ALL TOOLCHAIN WORKSHOP AMI'S"
    GET_ALL_AMIS ${TEMPLATE_LIST[@]}
  else
    GET_OLDER_DUPLICATE_AMIS ${1,,} ${TEMPLATE_LIST[@]}
  fi
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
