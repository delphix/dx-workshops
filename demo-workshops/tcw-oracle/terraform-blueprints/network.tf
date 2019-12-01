module "aws_vpc" {
  source  = "./modules/vpc"
  project = var.project
}

module "aws_subnet" {
  source            = "./modules/subnet"
  project           = var.project
  how_many          = var.how_many
  vpc_id            = module.aws_vpc.id
  availability_zone = var.availability_zone
}

module "aws_security_group" {
  source                             = "./modules/firewall"
  project                            = var.project
  how_many                           = var.how_many
  vpc_id                             = module.aws_vpc.id
  name                               = "${var.project}-${random_id.rval.hex}"
  addtl_firewall_ingress_cidr_blocks = var.addtl_firewall_ingress_cidr_blocks
  dev_mode                           = var.dev_mode
}

