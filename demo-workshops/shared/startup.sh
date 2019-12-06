#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
source $(dirname "${BASH_SOURCE[0]}")/library.sh
ARG=${1}
shift

case ${ARG} in
validate)
  exec /bin/validate.sh "$@"
  ;;
build)
  exec /bin/packer_build.sh "$@"
  ;;
deploy)
  exec /bin/terraform.sh apply "$@"
  ;;
redeploy)
  exec /bin/terraform.sh redeploy "$@"
  ;;
plan)
  exec /bin/terraform.sh plan "$@"
  ;;
teardown)
  exec /bin/terraform.sh destroy "$@"
  ;;
show)
  exec /bin/terraform.sh show "$@"
  ;;
output)
  exec /bin/terraform.sh output "$@"
  ;;
env|environment)
  exec /bin/terraform.sh output environment "$@"
  ;;
cleanup)
  exec /bin/cleanup.sh "$@"
  ;;
start|up)
  exec /bin/updown.sh start "$@"
  ;;
stop|down)
  exec /bin/updown.sh stop "$@"
  ;;
firewall|fw)
  exec /bin/terraform.sh apply -target="module.aws_security_group.aws_security_group.jumpbox" -target="module.aws_security_group.aws_security_group.landshark" "$@"
  ;;
image)
  exec /bin/image.sh "$@"
  ;;
ready)
  exec /bin/ws_ready.sh
  ;;
lab_print)
  LAB_PRINT "$@"
  ;;
""|*)
  help
  ;;
esac