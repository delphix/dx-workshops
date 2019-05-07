def packerBuild() {
  sh (
    script: """#!/bin/bash
      { set +x; } 2>/dev/null
      source ${env.IMAGEBUILDER_LIB}
      source /var/lib/jenkins/.packer_templates.env
      cd  \$(RETURN_DIRECTORY ${env.STAGE_NAME})
      echo "Changing instance_type to t2.2xlarge (37 cents/hour/machine) to minimize build time."
      sed -i -e 's|\\(.*instance_type\\"\\: \\"\\)\\(.*\\)\\(\\",\\)|\\1t2.2xlarge\\3|' ${env.STAGE_NAME}
      ln -s ${env.WORKSPACE}/ansible/roles roles
      PACKER_BUILD ${env.STAGE_NAME}
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
    echo "Changing all instance_types to t2.2xlarge (37 cents/hour/machine) to minimize build time."
    for each in `grep -R instance_type modules | awk -F: '{print \$1}'`
    do
      echo Updating instance_type in \$each
      sed -i -e \'s|\\(.*instance_type = \\"\\)\\(.*\\)\\(\\"\\)|\\1t2.2xlarge\\3|\' \${each}
    done

    cp /var/lib/jenkins/terraform_dev_backend.tf .
    terraform init
    terraform apply -var-file=${env.TF_VARS} --auto-approve -var "staged=${staged}"
  """
}

def terraformDestroy(){
  sh """#!/bin/bash
    { set -x; } 2>/dev/null
    terraform destroy -var-file=${env.TF_VARS} --auto-approve || true
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
