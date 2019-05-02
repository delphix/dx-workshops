variable "project" {
  description = "Short name of project"
}

variable "default_tags" {
  default = {}
}

variable "how_many" {
  default = 1
}

variable "vpc_id" {
  description = "the id of the VPC to link to"
}

variable "name" {
  description = "unique name for the sg"
}

variable "addtl_firewall_ingress_cidr_blocks" {
  description = "Specify any additional cidr block to allow inbound:any, i.e. '0.0.0.0/0' to allow all"
  //leave empty to allow external access only from the IP of the machine executing terraform
  default = []
}