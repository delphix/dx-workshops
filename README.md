# Various Pack Templates for images ready for use with the Delphix Dynamic Data Platform <!-- omit in toc -->

- [Description](#description)
- [Workshops](#workshops)
  - [Toolchain Workshop](#toolchain-workshop)
- [Base Templates](#base-templates)
  - [Usage](#usage)
    - [Configuring](#configuring)
    - [User Variables](#user-variables)
    - [Building](#building)
- [Links](#links)
- [Contribute](#contribute)
  - [Workflow](#workflow)
- [Reporting Issues](#reporting-issues)
- [Statement of Support](#statement-of-support)
- [License](#license)

## Description

This repo consists of some standard OS + Database templates configured ready to use with Delphix (located in the base-templates folder):

- delphix-centos6.9-oracle11.2.0.4.json
- delphix-centos7-ansible-base.json
- delphix-centos7-oracle-11.2.0.4.json
- delphix-centos7-oracle-12.2.0.1.json
- delphix-centos7-postrges-9.6.json
- delphix-ubuntu-bionic-guacamole.json

There are also some bespoke Packer templates used for specific scenarios:

- delphix-centos7-daf-app.json
- delphix-centos7-kitchen_sink.json
- delphix-centos7-tooling-base.json
- delphix-tcw-jumpbox.json
- delphix-tcw-oracle12-source.json
- delphix-tcw-oracle12-target.json
- delphix-tcw-tooling-oracle.json
- delphix-toolchain-dafdb-pgsql-source.json

The Oracle templates requires access to the Oracle binaries placed in an AWS s3 bucket.

## Workshops

### Toolchain Workshop

Please see the documentation for the [Toolchain Workshop](demo-workshops/tcw/docs/building.md)

## Base Templates

1. Clone this repository
2. Navigate into the cloned directory in a terminal
3. Copy the .example.env to .environment.env

These templates depend on Packer and Ansible existing in your path. If you are running a Mac, then the easiest way is to install via homebrew.
After cloning this repo, install the required Ansible dependencies.

```bash
brew install ansible packer git
git clone https://github.com/delphix/packer-templates
cd base-templates
cp .example.env .example.env
ansible-galaxy install -r roles/X_requirements.yml
```

### Usage

#### Configuring

1. Edit the .environment.env file in the root directory of the cloned repo

#### User Variables

1. AWS_ACCESS_KEY_ID - The AWS_ACCESS_KEY_ID environment variable
2. AWS_SECRET_ACCESS_KEY - The AWS_SECRET_ACCESS_KEY environment variable
3. AWS_REGION - The region packer will build the temporary infrastructure for the AMI
4. ORACLE_BINARIES_ROOT_URL - The URL where the AWS instance can retrieve the Oracle binaries during the packer build (For Oracle templates, only)
5. AWS_VPC_ID - The VPC ID from the region that packer will use
6. AWS_SUBNET_ID - The subnet ID from the VPC that packer will use

The below values are arbitrary, and only for tagging resources

1. AWS_EXPIRATION - The date this AMI is expired, i.e. "2037-07-01" or "never"
2. AWS_OWNER - The name of the person who owns the AMI, i.e. "Adam Bowen"
3. AWS_PROJECT - The name of the project that the AMI belongs, or came, from
4. AWS_COSTCENTER - The name of the cost center, if applicable

#### Building

1. source the .environment.env file
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

## Links

Include useful links to references or more advanced guides.

- [Packer Intro](https://www.packer.io/intro)
- [Building AMIs with Packer](https://www.packer.io/intro/getting-started/build-image.html)

## Contribute

All contributors are required to sign the Delphix Contributor Agreement prior to contributing code to an open source
repository. This process is handled automatically by [cla-assistant](https://cla-assistant.io/). Simply open a pull
request and a bot will automatically check to see if you have signed the latest agreement. If not, you will be prompted
to do so as part of the pull request process.

Please note that this project is released with a
[Contributor Code of Conduct](https://delphix.github.io/code-of-conduct.html). By participating in this project you
agree to abide by its terms.

### Workflow

1. Fork the project.
2. Make your bug fix or new feature.
3. Add tests for your code.
4. Send a pull request.

Contributions must be signed as `User Name <user@email.com>`. Make sure to [set up Git with user name and email address](https://git-scm.com/book/en/v2/Getting-Started-First-Time-Git-Setup). Bug fixes should branch from the current stable branch. New feature should be based on the `master` branch.

## Reporting Issues

Issues should be reported in the repo's issue tab.

## Statement of Support

This software is provided as-is, without warranty of any kind or commercial support through Delphix. See the associated
license for additional details. Questions, issues, feature requests, and contributions should be directed to the
community as outlined in the [Delphix Community Guidelines](https://delphix.github.io/community-guidelines.html).

## License

This is code is licensed under the Apache License 2.0. Full license is available [here](./LICENSE).
