module "delphix_target" {
  source                      = "./modules/delphix_target"
  project                     = var.project
  key_name                    = var.key_name
  ami_name                    = "delphix-tcw-target-oracle12-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  env_name                    = "TARGET"
  last_octet                  = "30"
  associate_public_ip_address = var.dev_mode
}

output "OracleTarget" {
  value = formatlist(
    "Oracle Target - Public IP: %s Private IP: %s",
    module.delphix_target.public_ip,
    module.delphix_target.private_ip,
  )
}

