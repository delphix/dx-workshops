module "delphix_source" {
  source                      = "./modules/delphix_target"
  project                     = var.project
  key_name                    = var.key_name
  ami_name                    = "delphix-tcw-source-oracle12-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  env_name                    = "SOURCE"
  last_octet                  = "20"
  associate_public_ip_address = var.dev_mode
}

output "OracleSource" {
  value = formatlist(
    "Oracle Source - Public IP: %s Private IP: %s",
    module.delphix_source.public_ip,
    module.delphix_source.private_ip,
  )
}

