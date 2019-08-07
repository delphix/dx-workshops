# FROM hashicorp/packer:light
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#

FROM golang:1.12-alpine
LABEL  maintainer="The Delphix Team"

ARG WORKSHOP_PATH
ARG DEMO_NAME

WORKDIR /build 

ENV DEMO_PATH=${WORKSHOP_PATH}/${DEMO_NAME}
ENV PACKER_VERSION=1.4.2
ENV ANSIBLE_VERSION=2.8.2
ENV ANSIBLE_ROLES_PATH=/.ansible
ENV TERRAFORM_VERSION=0.12.5
ENV AWSCLI_VERSION=1.16.201
ENV BOTO_VER=2.49.0
ENV BOTOCORE_VER=1.12.191
ENV BOTO3_VER=1.9.191
#DOCUTILS 0.15 breaks awscli
ENV DOCUTILS_VER=0.14
ENV FUTURES_VER=3.3.0
ENV JMESPATH_VER=0.9.4
ENV RSA_VER=3.4.2
ENV S3TRANSFER_VER=0.2.1
ENV URLLIB3_VER=1.25.3

RUN \
    apk add --no-cache \
    bash \
    curl \
    gcc \
    git \
    jq \
    libc-dev \
    make \
    openssh-client \
    python \
    python-dev \
    py-boto \
    py-dateutil \
    py-httplib2 \
    py-jinja2 \
    py-netaddr \
    py-paramiko \
    py-pip \
    py-setuptools \
    py-yaml \
    tar

#Packer
RUN curl -SL https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_SHA256SUMS --output packer_${PACKER_VERSION}_SHA256SUMS \
    && curl -SL https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip --output packer_${PACKER_VERSION}_linux_amd64.zip \
    && sed -i '/.*linux_amd64.zip/!d' packer_${PACKER_VERSION}_SHA256SUMS \
    && sha256sum -cs packer_${PACKER_VERSION}_SHA256SUMS \
    && unzip packer_${PACKER_VERSION}_linux_amd64.zip -d /bin \
    && rm -f packer_${PACKER_VERSION}_linux_amd64.zip packer_${PACKER_VERSION}_SHA256SUMS

#Terraform
RUN curl -SL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS --output terraform_${TERRAFORM_VERSION}_SHA256SUMS \
    && curl -SL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip --output terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && sed -i '/.*linux_amd64.zip/!d' terraform_${TERRAFORM_VERSION}_SHA256SUMS \
    && sha256sum -cs terraform_${TERRAFORM_VERSION}_SHA256SUMS \
    && unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /bin \
    && rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip terraform_${TERRAFORM_VERSION}_SHA256SUMS

#ansible, awscli, and the botos
RUN pip install -U pip \
    && pip install -U ansible==${ANSIBLE_VERSION} awscli==${AWSCLI_VERSION} boto==${BOTO_VER} \
    botocore==${BOTOCORE_VER} boto3==${BOTO3_VER} docutils==${DOCUTILS_VER} futures==${FUTURES_VER} \
    jmespath==${JMESPATH_VER} rsa==${RSA_VER} s3transfer==${S3TRANSFER_VER} urllib3==${URLLIB3_VER}

#install the ansible roles
COPY ${WORKDIR}/ansible/all_requirements.yml all_requirements.yml 
RUN ansible-galaxy install -r all_requirements.yml \
    && rm all_requirements.yml 

#compile the Golang binary used to gracefully shutdown the DDDP VDB's
COPY ${WORKDIR}/go/shutdown_dbs shutdown_dbs
RUN cd shutdown_dbs \
    && go get -d -v ./... \
    && env GOOS=linux GOARCH=amd64 go build -o /bin/shutdown_dbs \
    && rm -Rf shutdown_dbs

#place scripts in /bin
COPY ${DEMO_PATH}/image-builder/terraform.sh /bin/terraform.sh
RUN chmod +x /bin/terraform.sh

COPY ${DEMO_PATH}/image-builder/startup.sh /bin/startup.sh
RUN chmod +x /bin/startup.sh

COPY ${DEMO_PATH}/image-builder/library.sh /bin/library.sh
RUN chmod +x /bin/library.sh

COPY ${DEMO_PATH}/image-builder/packer_build.sh /bin/packer_build.sh
RUN chmod +x /bin/packer_build.sh

COPY ${DEMO_PATH}/image-builder/cleanup.sh /bin/cleanup.sh
RUN chmod +x /bin/cleanup.sh

COPY ${DEMO_PATH}/image-builder/validate.sh /bin/validate.sh
RUN chmod +x /bin/validate.sh

COPY ${DEMO_PATH}/image-builder/updown.sh /bin/updown.sh
RUN chmod +x /bin/updown.sh

COPY ${DEMO_PATH}/image-builder/image.sh /bin/image.sh
RUN chmod +x /bin/image.sh

COPY ${DEMO_PATH}/image-builder/ws_ready.sh /bin/ws_ready.sh
RUN chmod +x /bin/ws_ready.sh

RUN curl -SL https://github.com/mikefarah/yq/releases/download/2.4.0/yq_linux_amd64 --output /bin/yq \
    && chmod +x /bin/yq

ENTRYPOINT ["/bin/startup.sh"]