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
  if [[ "${1,,}" == "intermediate" || "${1,,}" == "--even-latest" || "${1,,}" == "all" ]]; then
    echo intermediate
    if [[ "${1,,}" == "all" ]]; then
      echo "WARNING: THIS WILL REMOVE ALL TOOLCHAIN WORKSHOP AMI'S"
      unset SUFFIX
    fi
    GET_ALL_AMIS $(GET_SYSTEMS)
  else
    SUFFIX=${1:-$SUFFIX}
    if [[ -n "${2}" || "${2,,}" == "--even-latest" ]]; then
      echo "suffix with EL"
      GET_ALL_AMIS $(GET_SYSTEMS)
    else
      echo "suffix without EL"
      GET_CLEANUP_LIST $(GET_SYSTEMS)
    fi
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
