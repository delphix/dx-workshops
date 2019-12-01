module "tooling" {
  source            = "./modules/tooling"
  project           = var.project
  key_name          = var.key_name
  how_many          = var.how_many
  subnet_id         = module.aws_subnet.id
  security_group_id = [module.aws_security_group.id]
  default_tags      = local.default_tags
  last_octet        = "6"
  ami_name          = "delphix-tcw-tooling-postgres11-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
}

output "TOOLING" {
  value = formatlist(
    "Tooling - Public IP: %s Private IP: %s",
    module.tooling.public_ip,
    module.tooling.private_ip,
  )
}

