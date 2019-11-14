module "virtualization_engine" {
  source                      = "./modules/delphix_engine"
  project                     = var.project
  key_name                    = var.key_name
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  last_octet                  = "10"
  staged                      = var.staged
  ami_name                    = var.staged == "false" ? "Delphix Engine ${var.delphix_engine_version}" : "delphix-tcw-virtualizationengine-${var.stage_name}-*"
  associate_public_ip_address = var.dev_mode
  engine_type                 = "VE"
}

output "DDDP-Virtualization" {
  value = formatlist(
    "Delphix Engine - Public IP: %s Private IP: %s    Access via browser @ http://%s",
    module.virtualization_engine.public_ip,
    module.virtualization_engine.private_ip,
    module.virtualization_engine.public_ip,
  )
}

