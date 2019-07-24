output "delphix-tcw-oracle12-source_ip" {
  value = module.delphix_source.public_ip
}

output "delphix-tcw-oracle12-source_id" {
  value = module.delphix_source.instance_id
}

output "delphix-tcw-oracle12-target_ip" {
  value = module.delphix_target.public_ip
}

output "delphix-tcw-oracle12-target_id" {
  value = module.delphix_target.instance_id
}

output "delphix-tcw-tooling-oracle_ip" {
  value = module.tooling.public_ip
}

output "delphix-tcw-tooling-oracle_id" {
  value = module.tooling.instance_id
}

output "delphix-tcw-virtualizationengine_ip" {
  value = module.virtualization_engine.public_ip
}

output "delphix-tcw-virtualizationengine_id" {
  value = module.virtualization_engine.instance_id
}

output "delphix-tcw-maskingengine_ip" {
  value = module.masking_engine.public_ip
}

output "delphix-tcw-maskingengine_id" {
  value = module.masking_engine.instance_id
}

output "delphix-tcw-jumpbox_ip" {
  value = module.guacamole.public_ip
}

output "delphix-tcw-jumpbox_id" {
  value = module.guacamole.instance_id
}

output "prodweb_ip" {
  value = module.prod_web_server.public_ip
}

output "prodweb_id" {
  value = module.prod_web_server.instance_id
}

output "devweb_ip" {
  value = module.dev_web_server.public_ip
}

output "devweb_id" {
  value = module.dev_web_server.instance_id
}

output "testweb_ip" {
  value = module.test_web_server.public_ip
}

output "testweb_id" {
  value = module.test_web_server.instance_id
}

