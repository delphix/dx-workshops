module "prod_web_server" {
  source                      = "./modules/web_server"
  project                     = var.project
  key_name                    = var.key_name
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  ami_name                    = "delphix-tcw-centos7-daf-app-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
  last_octet                  = "72"
  env_name                    = "prod"
  associate_public_ip_address = var.dev_mode
}

output "prod_public_ip" {
  value = module.prod_web_server.public_ip
}

output "prod_private_ip" {
  value = module.prod_web_server.public_ip
}

output "prod_db_host" {
  value = module.delphix_source.private_ip
}

output "ProdWebServer" {
  value = formatlist(
    "Prod Web Server - Public IP: %s Private IP: %s",
    module.prod_web_server.public_ip,
    module.prod_web_server.private_ip,
  )
}

