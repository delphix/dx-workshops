# Packer Template for CentOS 6.9 with Oracle 11.2.0.4 installed and prepared for Delphix

#### Table of Contents
1.  [Description](#description)
2.  [Installation](#installation)
3.  [Usage](#usage)
4.  [Links](#links)
5.  [Contribute](#contribute)
6.  [Reporting Issues](#reporting-issues)
7.  [License](#license)

## <a id="description"></a>Description

These are working Packer templates that will create AWS images that consists of:
- CentOS 6.9 w/ Oracle 11.2.0.4
- CentOS 7 w/ Oracle 11.2.0.4
- CentOS 7 w/ Postgres 9.6
- Delphix prerequisites installed and configured for use.

There are also some demo-specific Packer Templates
- Ubunutu Bionic w/Guacamole
- CentOS 7 configured for the RDS demo
- CentOS 7 configured as a source for the Delphix Automation Framework Demo


The Oracle templates require the provisioned system to be able to access p13390677_112040_Linux-x86-64_1of7.zip and p13390677_112040_Linux-x86-64_2of7.zip via an unauthenticated http location (like an s3 bucket). 

## <a id="installation"></a>Installation

### <a id="installation-via-docker"></a>Via Docker (the easiest) ###
1. Clone this repository
2. Navigate into the cloned directory
3. Copy the .example.docker to .environment.env

```bash
git clone https://github.com/delphix/packer-templates
cd packer-templates
cp .example.docker .example.env
```

### <a id="installation-without-docker"></a>Without Docker (the second easiest) ###
1. Clone this repository
2. Navigate into the cloned directory
3. Copy the .example.env to .environment.env
4. Certain packer templates have additional variables. Make copies of those example files, as in step 3

This template depends on packer and ansible existing in your path. If you are running a Mac, then the easiest way is to install via homebrew.
After cloning this repo, install the required ansible dependencies.

```bash
git clone https://github.com/delphix/packer-templates
cd packer-templates
brew install ansible packer
cp .example.env .example.env
ansible-galaxy install -r roles/requirements.yml
```

## <a id="usage"></a>Usage
### Configuring
1. Edit the .environment.env file ###(HERE)### in the root directory of the cloned repo
### user variables
1. AWS_ACCESS_KEY_ID - The AWS_ACCESS_KEY_ID environment variable
2. AWS_SECRET_ACCESS_KEY - The AWS_SECRET_ACCESS_KEY environment variable
3. AWS_BUILD_REGION - The region packer will build the temporary infrastructure for the AMI
4. AWS_DESTINATION_REGIONS - A comma-delimited list of regions for packer to copy the AMI
5. ORACLE_BINARIES_ROOT_URL - The URL where the AWS instance can retrieve the Oracle binaries during the packer build
6. AWS_VPC_ID - The VPC ID from the region that packer will use
7. AWS_SUBNET_ID - The subnet ID from the VPC that packer will use
#The below values are arbitrary, and only for tagging resources
8. AWS_EXPIRATION - The date this AMI is expired, i.e. "2037-07-01" or "never"
9. AWS_OWNER - The name of the person who owns the AMI, i.e. "Adam Bowen"
10. AWS_PROJECT - The name of the project that the AMI belongs, or came, from
11. AWS_COSTCENTER - The name of the cost center, if applicable

### Building
#### Via Docker
1. (Optional) Pull the cloudsurgeon/rds_demo docker image
2. Run the container against the packer template you want to use.

```bash
docker pull cloudsurgeon/packer-ansible
docker run --env-file .environment.env -v $(pwd):/build -i -t cloudsurgeon/packer-ansible:latest build delphix-centos7-rds.json
```

#### Without Docker
1. source the .example.env file
2. run packer against the template you want to use:
```bash
source .example.env
packer build delphix-centOS6.9-oracle11.2.0.4.json
```




```bash
packer build delphix-centOS6.9-oracle11.2.0.4.json 
cent69-Oracle11204 output will be in this color.

==> cent69-Oracle11204: Force Deregister flag found, skipping prevalidating AMI Name
    cent69-Oracle11204: Found Image ID: ami-8b44f234
==> cent69-Oracle11204: Creating temporary keypair: packer_5ad55a87-66df-e148-9439-a7bd06aa04fb
==> cent69-Oracle11204: Creating temporary security group for this instance: packer_5ad55ad1-10d5-7eda-e077-9741925ce7e4
==> cent69-Oracle11204: Authorizing access to port 22 from 0.0.0.0/0 in the temporary security group...
==> cent69-Oracle11204: Launching a source AWS instance...
==> cent69-Oracle11204: Adding tags to source instance
    cent69-Oracle11204: Adding tag: "Name": "Packer Builder"
    cent69-Oracle11204: Instance ID: i-0434f9b2e3a90c0c3
==> cent69-Oracle11204: Waiting for instance (i-0434f9b2e3a90c0c3) to become ready...
==> cent69-Oracle11204: Waiting for SSH to become available...
```

Second Example:
```bash
source .dafdb-source
packer build delphix-toolchain-dafdb-source.json
```




```bash
packer build delphix-toolchain-dafdb-source.json
delphix-toolchain-dafdb-source output will be in this color.

==> delphix-toolchain-dafdb-source: Force Deregister flag found, skipping prevalidating AMI Name
    delphix-toolchain-dafdb-source: Found Image ID: ami-090c2433423df7c1b
==> delphix-toolchain-dafdb-source: Creating temporary keypair: packer_5c24fd34-75a9-e10e-afa8-5784ac498ae2
==> delphix-toolchain-dafdb-source: Creating temporary security group for this instance: packer_5c24fd36-f5cc-8180-51b5-1bc18409d591
==> delphix-toolchain-dafdb-source: Authorizing access to port 22 from 0.0.0.0/0 in the temporary security group...
==> delphix-toolchain-dafdb-source: Launching a source AWS instance...
==> delphix-toolchain-dafdb-source: Adding tags to source instance
    delphix-toolchain-dafdb-source: Adding tag: "dlpx:CostCenter": "305000 - Development Engineering"
```

## <a id="links"></a>Links

Include useful links to references or more advanced guides.
*   [Packer Intro](https://www.packer.io/intro)
*   [Building AMIs with Packer](https://www.packer.io/intro/getting-started/build-image.html)

## <a id="contribute"></a>Contribute

Please note that this project is released with a [Contributor Code of Conduct](./code-of-conduct.md). By participating in this project you agree to abide by its terms.

#### Workflow

1.  Fork the project.
2.  Make your bug fix or new feature.
3.  Add tests for your code.
4.  Send a pull request.

Contributions must be signed as `User Name <user@email.com>`. Make sure to [set up Git with user name and email address](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup). Bug fixes should branch from the current stable branch. New feature should be based on the `master` branch.

## <a id="reporting_issues"></a>Reporting Issues

Issues should be reported in the repo's issue tab.

## <a id="license"></a>License

This is code is licensed under the Apache License 2.0. Full license is available [here](./LICENSE).
