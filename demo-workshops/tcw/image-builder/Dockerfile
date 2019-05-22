# FROM hashicorp/packer:light
#
# Copyright (c) 2019 by Delphix. All rights reserved.
#

FROM golang:alpine
MAINTAINER "The Delphix Team"

ARG DEMO_PATH
ARG DEMO_NAME

WORKDIR /build 

ENV PACKER_VERSION=1.3.5
ENV PACKER_SHA256SUM=822fe76c2dfe699f187ef8c44537d10453a1545db620e40b345cf6991a690f7d
ENV ANSIBLE_VERSION=2.7.9
ENV ANSIBLE_ROLES_PATH=/.ansible
ENV TERRAFORM_VERSION=0.11.13
ENV TERRAFORM_SHA256SUM=6b8ce67647a59b2a3f70199c304abca0ddec0e49fd060944c26f666298e23418
ENV AWSCLI_VERSION=1.16.149

RUN \
    apk add --no-cache \
    bash \
    curl \
    gcc \
    git \
    jq \
    libc-dev \
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

RUN pip install -U pip \
    && pip install -U ansible==${ANSIBLE_VERSION} awscli==${AWSCLI_VERSION}

COPY ${DEMO_PATH}/ansible/all_requirements.yml all_requirements.yml 
RUN ansible-galaxy install -r all_requirements.yml \
    && rm all_requirements.yml 

COPY ${DEMO_PATH}/${DEMO_NAME}/image-builder/terraform.sh /bin/terraform.sh
RUN chmod +x /bin/terraform.sh

COPY ${DEMO_PATH}/${DEMO_NAME}/image-builder/startup.sh /bin/startup.sh
RUN chmod +x /bin/startup.sh

COPY ${DEMO_PATH}/${DEMO_NAME}/image-builder/library.sh /bin/library.sh
RUN chmod +x /bin/library.sh

COPY ${DEMO_PATH}/${DEMO_NAME}/image-builder/packer_build.sh /bin/packer_build.sh
RUN chmod +x /bin/packer_build.sh

COPY ${DEMO_PATH}/${DEMO_NAME}/image-builder/cleanup.sh /bin/cleanup.sh
RUN chmod +x /bin/cleanup.sh
#this last line is so we can run packer/ansible as another user than root
RUN chmod 777 /.ansible

ENTRYPOINT ["/bin/startup.sh"]