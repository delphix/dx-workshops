#!/usr/bin/env bash
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#
STARTTIME=$(date +%s)
NOW=$(date +"%m-%d-%Y %T")
WORKDIR=$(pwd)

WORKSHOP_PATH="demo-workshops"
DEMO_NAME="tcw"
DEMO_PATH="${WORKDIR}/${WORKSHOP_PATH}/${DEMO_NAME}"
PACKER_TEMPLATES="${WORKDIR}/packer-templates"
CERT="${WORKDIR}/certs/ansible"
TERRAFORM_BLUEPRINTS="${DEMO_PATH}/terraform-blueprints"
GODIR="${WORKDIR}/go"

function cleanup() {
  ERROR
}

function GET_DEFAULT_AMI_SUFFIX {
  yq r ${DEMO_PATH}/workshop.yaml ami_suffix
}

SUFFIX=$(GET_DEFAULT_AMI_SUFFIX)

function ENVCHECK() {
	[[ -z "${S3_AWS_ACCESS_KEY_ID}" ]] && export S3_AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
	[[ -z "${S3_AWS_SECRET_ACCESS_KEY}" ]] && export S3_AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
  env | grep _KEY
}

function ENDTIME() {
	ENDTIME=$(date +%s)
	echo "It took $(($ENDTIME - $STARTTIME)) seconds to complete ${0}"
}

function ERROR() {
	ENDTIME
	[[ -f ${WORKDIR}/WAIT.log ]] && mv ${WORKDIR}/WAIT.log ${WORKDIR}/ERROR.log
	exit 1
}

function READY() {
	ENDTIME
	echo "Script finished Successfully"
	mv ${WORKDIR}/WAIT.log ${WORKDIR}/READY.log
}

function JOB_WAIT() {
	for job in `jobs -p`
	do
		wait $job || let "FAIL+=1"
	done
	[[ -n "${FAIL}" ]] && ERROR
  true #exits function cleanly, if no error condition
}

function CERT_TEST() {
	#Check to see if the automation certificates exist, if not prompt to create
	if [[ -f ${CERT}.pub && -f ${CERT} ]] ; then
		echo "Certifcate: ${CERT}"
	elif [[ ! -f ${CERT}.pub && -f ${CERT} ]]; then
    until [[ "${CREATEPUB}" == "y" || "${CREATEPUB}" == "n" ]]; do
      read -p "${CERT}.pub is missing. Should I create it? (y) " CREATEPUB
      CREATEPUB=${CREATEPUB:-y}
    done
    if [[ "${CREATEPUB}" == "y" ]]; then
      ssh-keygen -y -f ${CERT} > ${CERT}.pub
    else
      echo "set the CERT variable to the private certificate you would like to use for automation"
      echo "The public certificate should be found in the same directory"
      echo "This should be a key exclusive to run in this environment"
      exit 1
    fi
  else
    until [[ "${CREATECERT}" == "y" || "${CREATECERT}" == "n" ]]; do
      read -p "${CERT} is missing. Should I create it? (y) " CREATECERT
      CREATECERT=${CREATECERT:-y}
    done
    if [[ "${CREATECERT}" == "y" ]]; then
      ssh-keygen -b 2048 -t rsa -f certs/ansible -q -N ""
    else
      echo "set the CERT variable to the private certificate you would like to use for automation"
      echo "The public certificate should be found in the same directory"
      echo "This should be a key exclusive to run in this environment"
      exit 1
    fi
	fi
	export ANSIBLE_PUB_KEY=$(cat ${CERT}.pub)
}

function AMI_INFO() {
    for each in "$@"; do
        # query for an existing AMI with the name and md5sum number, and store that information in a file
        aws ec2 --region ${AWS_REGION} describe-images --owner self --filters "Name=name,Values=${each}-*" "Name=tag:md5sum,Values=$(GET_MD5SUM ${each})" --query 'sort_by(Images, &CreationDate)[-1]' > /tmp/${each}_info.json
    done
    JOB_WAIT
}

function AMI_EXISTS() {
	# if a an ImageId is present, then the ami exists
    AMI=$(jq -r '.ImageId' /tmp/${1}_info.json)
	[[ -n $AMI && $AMI != 'null' ]] 
}

function AMI_OLDER_THAN_SOURCE() {
	# Check the source tag of the AMI and compare the CreationDate of the AMI and it's Source, if applicable
	SOURCE_AMI="$(cat /tmp/${1%.json}_info.json | jq -r '.Tags[]|select(.Key=="source").Value')" 
    if [[ -f /tmp/${SOURCE_AMI%.json}_info.json ]]; then
        echo comparing ${1%.json}: $(cat /tmp/${1%.json}_info.json | jq -r '.CreationDate') to ${SOURCE_AMI%.json}: $(cat /tmp/${SOURCE_AMI%.json}_info.json | jq -r '.CreationDate')
        [[ "$(cat /tmp/${1%.json}_info.json | jq -r '.CreationDate')" < "$(cat /tmp/${SOURCE_AMI%.json}_info.json | jq -r '.CreationDate')" ]]
    else
        false
    fi
}

function NEED_TO_BUILD_AMI() {
  if AMI_EXISTS $1; then
      # If an AMI with a valid name and md5sum exists
      if AMI_OLDER_THAN_SOURCE $1; then
          # Check to see if it is older than it's source AMI. If so, build it
          echo ${1} exists, but is older than source. Rebuilding.
          true
      else
          echo Skipping ${1}: a valid ami exists
          false
      fi
  else
      # Otherwise, build it
      echo Building ${1}
      true
  fi
}

function PACKER_BUILD() {
    SYSTEMS=$(GET_BATCH_SYSTEMS ${1})
    AMI_INFO ${SYSTEMS}
    for each in ${SYSTEMS}; do
      if NEED_TO_BUILD_AMI $each; then
        #drop this little file to alert our build server know a new ami was built
        touch ${WORKDIR}/change.ignore
        cd $PACKER_TEMPLATES
        packer build -var "md5sum=$(GET_MD5SUM $each)" -var "ami_name_prefix=$(RETURN_AMI_NAMES ${each})" $(GET_SYSTEM_PACKER_TEMPLATE ${each}) &
        cd - &>>/dev/null
      fi
    done
    JOB_WAIT
    AMI_INFO SYSTEMS
}

function RETURN_DIRECTORY() {
	#This function is a placeholder until I figure how I want to reorg folder structure
	dirname $(find $WORKDIR -name $1)
}

function GET_CLEANUP_LIST() {
  local LIST
  local FILTER
  for SYSTEM in "$@"; do
    unset LIST
    local AMI_NAME=$(RETURN_AMI_NAMES ${SYSTEM})
    local AMI_LIST=$(aws ec2 --region ${AWS_REGION} describe-images --owner self --filters "Name=name,Values=${AMI_NAME}-*" )
    # query for an existing AMI with the name and md5sum number, and store that information in a file
    echo "Will deregister the following AMI's for ${AMI_NAME}-*:"
    if [[ "$(yq r ${DEMO_PATH}/workshop.yaml ami_suffix)" == "${SUFFIX}" ]]; then
      FILTER=".Images[]|select((.Tags[]|select(.Key==\"md5sum\")|.Value) !=\"$(GET_MD5SUM ${SYSTEM})\" )|.ImageId" 
      for each in $( echo "${AMI_LIST}"| jq -r "${FILTER}" ); do
        LIST+=( ${each} )
        CLEANUP_LIST+=( ${each} )
      done
      FILTER="[.Images|sort_by(.CreationDate)|.[]|select((.Tags[]|select(.Key==\"md5sum\")|.Value) ==\"$(GET_MD5SUM ${SYSTEM})\")| .ImageId][0:-1]|.[]"
    else
      FILTER="[.Images|sort_by(.CreationDate)|.[]| .ImageId][0:-1]|.[]"
    fi
    for each in $(echo ${AMI_LIST}| jq -r "${FILTER}"); do
      LIST+=( ${each} )
      CLEANUP_LIST+=( ${each} )
    done 
    printf "%s\n" "${LIST[@]}" | sort -u
  done
}

function CLEANUP_AMIS() {
  for each in $(printf "%s\n" "${CLEANUP_LIST[@]}" | sort -u); do
    echo "Deregistering ${each}"
    aws ec2 --region ${AWS_REGION} deregister-image --image-id ${each}
  done
}

#Grab all AMI's that match the name prefix
function GET_ALL_AMIS() {
  for each in "$@"; do
    AMI_NAME=$(RETURN_AMI_NAMES ${each})
    echo "Will deregister the following AMI's for ${AMI_NAME}:"
    for ami in $(aws ec2 --region ${AWS_REGION} describe-images --owner self --filters "Name=name,Values=${AMI_NAME}-*" --query 'Images[*].ImageId' --output text); do
      echo -e "\t${ami}"
      CLEANUP_LIST+=( ${ami} )
    done
  done
}

function SHUTDOWN_VDBS(){
  cd ${TERRAFORM_BLUEPRINTS}
  terraform refresh
  DE=$(terraform output -json delphix-tcw-virtualizationengine_ip | jq -r '.[]')
  if [[ -n $DE ]] ; then
    sed -e 's|ddp_hostname.*|ddp_hostname = '${DE}'|' \
      -e 's|password.*|password = '${DELPHIX_ADMIN_PASSWORD}'|' \
      -e 's|username.*|username = delphix_admin|' \
      ${GODIR}/shutdown_dbs/example_conf.txt > /tmp/shutdown_conf.txt
    shutdown_dbs -c /tmp/shutdown_conf.txt
  fi
}

function BINARY_BUILD() {
  cd $GODIR
  for each in `ls`; do
    cd ${GODIR}/${each}
    make build
  done
}

function GET_BATCHES() {
  yq r ${DEMO_PATH}/workshop.yaml --tojson | jq -r  '[.amis[]| select(.[]?.packer.name != null) | .[].packer.batch]|unique|.[]'
}

function GET_BATCH_SYSTEMS() {
  [[ -z ${1} ]] && echo "GET_BATCH_SYSTEMS requires a batch number" && ERROR
  FILTER=".amis[]| select(.[]?.packer.batch == ${1}) | to_entries[].key"
  yq r ${DEMO_PATH}/workshop.yaml --tojson | jq -r  "${FILTER}"
}

function GET_SYSTEM_PACKER_TEMPLATE() {
  [[ -z ${1} ]] && echo "GET_SYSTEM_PACKER_TEMPLATE requires a system name" && ERROR
  yq r ${DEMO_PATH}/workshop.yaml --tojson | jq -r  ".amis[]|select(.[]?.packer.name != null)|to_entries[]|select(.key == \"${1}\")|.value.packer.name"
}

function RETURN_TEMPLATES_FROM_BATCH() {
  unset TEMPLATE_LIST
  [[ -z ${1} ]] && echo RETURN_TEMPLATES_FROM_BATCH requires the BATCH number && ERROR
  for SYSTEM in `GET_BATCH_SYSTEMS ${1}`; do
    TEMPLATE_LIST="${TEMPLATE_LIST} $(GET_SYSTEM_PACKER_TEMPLATE ${SYSTEM})"
  done
  echo ${TEMPLATE_LIST}
}

function GET_SYSTEMS() {
  FILTER=".amis[] | to_entries[].key"
  yq r ${DEMO_PATH}/workshop.yaml --tojson | jq -r "${FILTER}"
}

function GET_MD5SUM() {
  [[ -z ${1} ]] && echo GET_MD5SUM requires the SYSTEM && ERROR 
  FILTER=".amis[]|select(.[]?.packer.name != null)|to_entries[]|select(.key == \"${1}\")|.value.packer.md5sum"
  yq r ${DEMO_PATH}/workshop.yaml --tojson | jq -r  "${FILTER}"
}

function RETURN_AMI_NAMES {
  local AMIS
  if [[ -n ${SUFFIX} ]]; then
    for SYSTEM in "${@}"; do
      AMIS+=( ${SYSTEM}-${SUFFIX} )
    done
  else
    AMIS="${@}"
  fi
  echo $AMIS
}