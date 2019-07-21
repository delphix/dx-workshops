module "guacamole" {
  source            = "./modules/guacamole"
  project           = var.project
  key_name          = var.key_name
  how_many          = var.how_many
  subnet_id         = module.aws_subnet.id
  security_group_id = [module.aws_security_group.id, module.aws_security_group.jumpbox]
  default_tags      = local.default_tags
  last_octet        = "5"
  ami_name          = "delphix-tcw-jumpbox-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
}

output "GUACAMOLE" {
  value = formatlist(
    "\nJumpbox - Public IP: %s Private IP: %s\n    Access via http://%s:8080/labs\n    Username: delphix",
    module.guacamole.public_ip,
    module.guacamole.private_ip,
    module.guacamole.public_ip,
  )
}

