#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
set -e
WORKDIR=$(pwd)
ARG=${1}
shift

function help() {
  echo """Usage: docker-compose run tcw <command> [args]

  Available commands are listed below.

  build              Build the images from the packer-templates folder
  cleanup            Cleanup all but the latest AMI's for this workshop
  deploy             Builds or changes Terraform-managed infrastructure
  plan               Generate and show an Terrafrom execution plan
  teardown           Destroy Terraform-managed infrastructure
  show               Print the terraform state
  env|environment    Print the jumpbox information
  """

}

case ${ARG} in
build)
  exec /bin/packer_build.sh "$@"
  ;;
deploy)
  exec /bin/terraform.sh apply "$@"
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
env|environment)
  exec /bin/terraform.sh output environment "$@"
  ;;
cleanup)
  exec /bin/cleanup.sh "$@"
  ;;
*)
  help
  ;;
esac