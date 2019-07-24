def packerBuild() {
  sh (
    """#!/bin/bash
      { set +x; } 2>/dev/null
      docker-compose run tcw build
      sudo /bin/chownership
    """
  )
}

def amiify() {
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    docker-compose run tcw image staged
  """
}

def environmentTest(){
  sh """#!/bin/bash
    docker-compose run tcw ready
  """
}

def terraformBuild(staged){
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    #copy in builderserver specific variables
    cp /var/lib/jenkins/build.auto.tfvars .
    CURRENT_UID=\$(id -u):\$(id -g) docker-compose run tcw deploy -auto-approve -var "staged=${staged}"
  """
}

def terraformDestroy(){
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    CURRENT_UID=\$(id -u):\$(id -g) docker-compose run tcw teardown -auto-approve || true
  """
}

def ami_checker(staged){
  sh (
    script: """#!/bin/bash
      { set +x; } 2>/dev/null
      source ${env.IMAGEBUILDER_LIB}
      source /var/lib/jenkins/.packer_templates.env
      AMI_INFO ${env.STAGE_NAME}
      NEED_TO_BUILD_AMI ${env.STAGE_NAME}
    """,
    returnStatus: true
    )
}
