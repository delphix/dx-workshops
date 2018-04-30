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

This is a working Packer template that will create an AWS image that consists of:
- CentOS 6.9
- Oracle 11.2.0.4
- Delphix prerequisites installed and configured for use.

This Packer template requires the provisioned system to be able to access p13390677_112040_Linux-x86-64_1of7.zip and p13390677_112040_Linux-x86-64_2of7.zip via an unauthenticated http location. 

## <a id="installation"></a>Installation

This template depends on packer and ansible existing in your path. If you are running a Mac, then the easiest way is to install via homebrew.
After cloning this repo, install the required ansible dependencies.

```bash
git clone https://gitlab.delphix.com/cto-research/packer-CentOS6.9-Oracle11.2.0.4
packer-CentOS6.9-Oracle11.2.0.4
brew install ansible packer
ansible-galaxy install -r requirements.yml
```

## <a id="usage"></a>Usage

1. Edit the user variables in .example.env to suit your environment
2. source the .example.env file
3. Build the ami.

```bash
#edit the .example.env file
vi .example.env
#...
source .example.env
packer build delphix-centOS6.9-oracle11.2.0.4.json
```

### user variables
1. AWS_ACCESS_KEY_ID - The AWS_ACCESS_KEY_ID environment variable
2. AWS_SECRET_ACCESS_KEY - The AWS_SECRET_ACCESS_KEY environment variable
3. AWS_BUILD_REGION - The region packer will build the temporary infrastructure for the AMI
4. AWS_DESTINATION_REGIONS - A comma-delimited list of regions for packer to copy the AMI
5. ORACLE_BINARIES_ROOT_URL - The URL where the AWS instance can retrieve the Oracle binaries during the packer build
6. AWS_VPC_ID - The VPC ID from the region that packer will use
7. AWS_SUBNET_ID - The subnet ID from the VPC that packer will use
8. AWS_EXPIRATION - The date this AMI is expired, i.e. "2037-07-01" or "never"
9. AWS_OWNER - The name of the person who owns the AMI, i.e. "Adam Bowen"
10. AWS_PROJECT - The name of the project that the AMI belongs, or came, from


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
