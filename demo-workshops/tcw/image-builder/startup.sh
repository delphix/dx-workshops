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
  validate           runs a few checks on the prereqs
  env|environment    Print the jumpbox information
  start [wait]       Start the stopped EC2 instances in an deployed environment (specifying wait will wait)
  stop [wait]        Stop the running EC2 instances in an deployed environment (specifying wait will wait)
  
  ex.
    docker-compose run tcw validate
    docker-compose run tcw build
    docker-compose run tcw deploy
    docker-compose run tcw env
    docker-compose run tcw stop wait
  
  """


}

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
start|up)
  exec /bin/updown.sh start "$@"
  ;;
stop|down)
  exec /bin/updown.sh stop "$@"
  ;;
*)
  help
  ;;
esac