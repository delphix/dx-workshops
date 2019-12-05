def packerBuild(workshop) {
  sh (
    script: """#!/bin/bash
      { set +x -e; } 2>/dev/null
      docker-compose run --rm ${workshop} build
    """
  )
}

def amiify(workshop) {
  sh """#!/bin/bash
    { set -x -e; } 2>/dev/null
    docker-compose run --rm ${workshop} image staged
  """
}

def environmentReady(workshop){
  sh """#!/bin/bash
    { set -x -e; } 2>/dev/null
    docker-compose run --rm ${workshop} ready
    #Adding in the below, to skip the reconfiguration since our environment builds are consistent
    JUMP=\$(docker-compose run ${workshop} output -json delphix-tcw-jumpbox_ip | tail -2 | head -1 | jq -r '.[]')
    ssh -i ${env.ANSIBLE_CERT} -o StrictHostKeyChecking=no ubuntu@\${JUMP} 'touch UPDOWN'
  """
}

def environmentTest(workshop){
  sh """#!/bin/bash
    { set -x -e; } 2>/dev/null
    JUMP=\$(docker-compose run --rm ${workshop} output -json delphix-tcw-jumpbox_ip | tail -2 | head -1 | jq -r '.[]')
    ssh -i ${env.ANSIBLE_CERT} -o StrictHostKeyChecking=no ubuntu@\${JUMP} 'cd tests;./tests.sh'
  """
}

def terraformBuild(staged, workshop){
  sh """#!/bin/bash
    { set -x -e; } 2>/dev/null
    #copy in builderserver specific variables
    cp /var/lib/jenkins/build.auto.tfvars .
    CURRENT_UID=\$(id -u):\$(id -g) docker-compose run --rm ${workshop} deploy -auto-approve -var "staged=${staged}"
  """
}

def terraformDestroy(workshop){
  sh """#!/bin/bash
    { set -x -e; } 2>/dev/null
    CURRENT_UID=\$(id -u):\$(id -g) docker-compose run --rm ${workshop} teardown -auto-approve || true
  """
}

def chownership() {
  sh (
   script: """#!/bin/bash
      { set +x -e; } 2>/dev/null
      sudo /bin/chownership
    """
  )
}