variable "access_key" {
}

variable "secret_key" {
}

variable "aws_region" {
  description = "The AWS region to create things in."
  default     = "us-west-2"
}

variable "availability_zone" {
  description = "Name of the availability zone to use"
  default     = "us-west-2a"
}

variable "key_name" {
  description = "Name of the SSH keypair to use in AWS."
}

variable "how_many" {
  default = 1
}

variable "expiration" {
  default = "24h"
}

variable "owner" {
  default = "073575380803"
}

variable "project" {
  description = "Short name of project"
}

variable "cost_center" {
  default = "305000 - Development Engineering"
}

data "aws_caller_identity" "current" {
}

locals {
  default_tags = {
    "dlpx:Project"    = var.project
    "dlpx:Expiration" = substr(timeadd(timestamp(), var.expiration), 0, 10)
    "dlpx:CostCenter" = var.cost_center
    "dlpx:Owner"      = data.aws_caller_identity.current.arn
    "UUID"            = random_id.rval.hex
    "STUDENT"         = ""
  }
}

variable "staged" {
  default = "false"
}

variable "stage_name" {
  default = "staged"
}

variable "delphix_engine_version" {
  //Currently only works with 5.3.2
  description = "The version of the Delphix AMI"
  default     = "5.3.2.*"
}

variable "addtl_firewall_ingress_cidr_blocks" {
  description = "Specify any additional cidr block to allow inbound:any, i.e. '0.0.0.0/0' to allow all (not recommended)"
  default     = []
}

variable "associate_public_ip_address" {
  description = "Associate public IP Address to every instance in the workshop"
  default     = "false"
}

