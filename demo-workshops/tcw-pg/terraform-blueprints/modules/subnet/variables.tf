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

variable "cidr_blocks" {
  description = "A list of CIDR blocks that can be used by the corresponding student ID's"
  default = {
  "0" = "10.0.0.0/24"
  "1" = "10.0.1.0/24"
  "2" = "10.0.2.0/24"
  "3" = "10.0.3.0/24"
  "4" = "10.0.4.0/24"
  "5" = "10.0.5.0/24"
  "6" = "10.0.6.0/24"
  "7" = "10.0.7.0/24"
  "8" = "10.0.8.0/24"
  "9" = "10.0.9.0/24"
  "10" = "10.0.10.0/24"
  "11" = "10.0.11.0/24"
  "12" = "10.0.12.0/24"
  "13" = "10.0.13.0/24"
  "14" = "10.0.14.0/24"
  "15" = "10.0.15.0/24"
  "16" = "10.0.16.0/24"
  "17" = "10.0.17.0/24"
  "18" = "10.0.18.0/24"
  "19" = "10.0.19.0/24"
  "20" = "10.0.20.0/24"
  }
}

variable "availability_zone" {
  description = "Name of the availability zone to use"
}