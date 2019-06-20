#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
trap "cleanup" SIGINT

{
    cd ${TERRAFORM_BLUEPRINTS}
    JUMP=$(terraform output delphix-tcw-jumpbox_ip)
    [[ -z $JUMP ]] && exit 1
    until ssh -i ${CERT} -o "StrictHostKeyChecking=no" ubuntu@${JUMP} 'ls ~/Desktop'|egrep "READY|ERROR"
    do
        ssh -i ${CERT} -o "StrictHostKeyChecking=no" ubuntu@${JUMP} 'tail -5 ~/Desktop/WAIT' || true
        sleep 10
    done
    if ! ssh -i ${CERT} -o "StrictHostKeyChecking=no" ubuntu@${JUMP} 'tail -5 ~/Desktop/READY' ; then
        ssh -i ${CERT} -o StrictHostKeyChecking=no ubuntu@${JUMP} 'cat ~/Desktop/ERROR'
        exit 1
    fi
}