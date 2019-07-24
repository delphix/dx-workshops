module "masking_engine" {
  source                      = "./modules/delphix_engine"
  project                     = var.project
  key_name                    = var.key_name
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  last_octet                  = "11"
  staged                      = var.staged
  ami_name                    = var.staged == "false" ? "Delphix Engine ${var.delphix_engine_version}" : "delphix-tcw-maskingengine-${var.stage_name}-*"
  associate_public_ip_address = var.associate_public_ip_address
  engine_type                 = "ME"
}

output "DDDP-Masking" {
  value = formatlist(
    "\nDelphix Engine - Public IP: %s Private IP: %s\n    Access via browser @ http://%s",
    module.masking_engine.public_ip,
    module.masking_engine.private_ip,
    module.masking_engine.public_ip,
  )
}

