def packerBuild() {
  sh (
    script: """#!/bin/bash
      { set +x; } 2>/dev/null
      docker-compose run tcw build
    """
  )
}

def amiify() {
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    source ${env.IMAGEBUILDER_LIB}
    . /var/lib/jenkins/.packer_templates.env
    instance_id=\$(terraform output ${env.STAGE_NAME}_id)
    AMI_NAME=${env.STAGE_NAME}
    ansible-playbook -i 'localhost,' ${env.WORKSPACE}/demo-workshops/ansible/ami_maker.yml -e "instance_id=\${instance_id}" -e "ami_name=\${AMI_NAME%.json}-staged" -e "commit=${env.GIT_COMMIT}"
  """
}

def environmentTest(){
  sh """#!/bin/bash
      { set -x; } 2>/dev/null
      GUAC=\$(terraform output delphix-tcw-jumpbox_ip)
      [[ -z \$GUAC ]] && exit 1
      until ssh -i ${env.ANSIBLE_CERT} -o "StrictHostKeyChecking=no" ubuntu@\${GUAC} 'ls ~/Desktop'|egrep "READY|ERROR"
      do
        ssh -i ${env.ANSIBLE_CERT} -o "StrictHostKeyChecking=no" ubuntu@\${GUAC} 'tail -5 ~/Desktop/WAIT' || true
        sleep 10
      done
      if ! ssh -i ${env.ANSIBLE_CERT} -o "StrictHostKeyChecking=no" ubuntu@\${GUAC} 'tail -5 ~/Desktop/READY' ; then
        ssh -i ${env.ANSIBLE_CERT} -o StrictHostKeyChecking=no ubuntu@\${GUAC} 'cat ~/Desktop/ERROR'
        exit 1
      fi 
  """
}

def terraformBuild(staged){
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    docker-compose run tcw deploy -auto-approve
  """
}

def terraformDestroy(){
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    docker-compose run tcw teardown -auto-approve || true
    sh "sudo rm -Rf ${env.TF_DIR}.terraform"
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
