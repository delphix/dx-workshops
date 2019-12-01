variable "subnet_id" {
  description = "The network for the delphix engine and target host"
  default = []
}

variable "ami_name" {
  description = "Name of the congigured target AMI."
}


variable "key_name" {
  description = "Name of the SSH keypair to use in AWS (do not include .pem extension)."
}

variable "how_many" {
  default = 1
}

variable "security_group_id" {
  description = "List of security group IDs"
  default = []
}

variable "default_tags" {
  default = {}
}

variable "project" {
  description = "Short name of project"
}

variable "env_name" {
  description = "The name of the environment in Delphix"
}

variable "last_octet" {
  description = "The last octet of the system's IP address"
}

variable "associate_public_ip_address" {
  description = "Associate public IP Address to every instance in the workshop"
  default = "false"
}