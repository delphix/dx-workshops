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

  image <name>       Shutdown running instances and create AMI's with name
  build              Build the images from the packer-templates folder
  cleanup            Cleanup all but the latest AMI's for this workshop
  deploy [args]      Builds or changes Terraform-managed infrastructure
  plan [args]        Generate and show a Terrafrom execution plan
  teardown [args]    Destroy Terraform-managed infrastructure
  redeploy [args]    Executes teardown and then deploy
  show               Print the terraform state
  validate           runs a few checks on the prereqs
  env|environment    Print the jumpbox information
  start [wait]       Start the stopped EC2 instances in a deployed environment
                     (specifying wait will wait)
  stop [wait]        Stop the running EC2 instances in a deployed environment
                     (specifying wait will wait)
  fw|firewall        Execute terraform apply against the firewall modules only (for firewall rule updates)
  
  ex.
    docker-compose run tcw validate
    docker-compose run tcw build
    docker-compose run tcw deploy
    docker-compose run tcw deploy --auto-approve
    docker-compose run tcw deploy -var \"staged=true\" -var \"stage_name=staged\"
    docker-compose run tcw env
    docker-compose run tcw stop wait
    docker-compose run tcw image staged
    docker-compose run tcw fw
  
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
*)
  help
  ;;
esac