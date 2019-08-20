module "dev_web_server" {
  source                      = "./modules/web_server"
  project                     = var.project
  key_name                    = var.key_name
  how_many                    = var.how_many
  subnet_id                   = module.aws_subnet.id
  security_group_id           = [module.aws_security_group.id]
  default_tags                = local.default_tags
  ami_name                    = "delphix-tcw-centos7-daf-app-${var.staged == "false" ? "unstaged" : var.stage_name}-*"
  last_octet                  = "71"
  env_name                    = "nonprod"
  associate_public_ip_address = var.dev_mode
}

output "dev_public_ip" {
  value = module.dev_web_server.public_ip
}

output "dev_private_ip" {
  value = module.dev_web_server.public_ip
}

output "DevWebServer" {
  value = formatlist(
    "Dev Web Server - Public IP: %s Private IP: %s",
    module.dev_web_server.public_ip,
    module.dev_web_server.private_ip,
  )
}

