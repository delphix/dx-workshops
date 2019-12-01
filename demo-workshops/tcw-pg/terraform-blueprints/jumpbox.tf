module "jumpbox" {
  source            = "./modules/jumpbox"
  project           = var.project
  key_name          = var.key_name
  how_many          = var.how_many
  subnet_id         = module.aws_subnet.id
  security_group_id = [module.aws_security_group.id, module.aws_security_group.jumpbox]
  default_tags      = local.default_tags
  last_octet        = "5"
  ami_name          = "delphix-tcw-jumpbox-postgres11-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
}

output "JUMPBOX" {
  value = formatlist(
    "Jumpbox - Public IP: %s Private IP: %s    Access via http://%s:8080/labs    Username: delphix",
    module.jumpbox.public_ip,
    module.jumpbox.private_ip,
    module.jumpbox.public_ip,
  )
}

