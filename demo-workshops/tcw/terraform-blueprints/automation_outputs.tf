output "delphix-tcw-oracle12-source_ip" {
  value = module.delphix_source.public_ip
}

output "delphix-tcw-oracle12-source_image" {
  value = module.delphix_source.instance_id
}

output "delphix-tcw-oracle12-target_ip" {
  value = module.delphix_target.public_ip
}

output "delphix-tcw-oracle12-target_image" {
  value = module.delphix_target.instance_id
}

output "delphix-tcw-tooling-oracle_ip" {
  value = module.tooling.public_ip
}

output "delphix-tcw-tooling-oracle_image" {
  value = module.tooling.instance_id
}

output "delphix-tcw-virtualizationengine_ip" {
  value = module.virtualization_engine.public_ip
}

output "delphix-tcw-virtualizationengine_image" {
  value = module.virtualization_engine.instance_id
}

output "delphix-tcw-maskingengine_ip" {
  value = module.masking_engine.public_ip
}

output "delphix-tcw-maskingengine_image" {
  value = module.masking_engine.instance_id
}

output "delphix-tcw-jumpbox_ip" {
  value = module.guacamole.public_ip
}

output "delphix-tcw-jumpbox_image" {
  value = module.guacamole.instance_id
}

output "delphix-tcw-centos7-daf-app_image" {
  value = module.prod_web_server.instance_id
}

output "delphix-tcw-oracle12-source_system" {
  value = module.delphix_source.instance_id
}

output "delphix-tcw-oracle12-target_system" {
  value = module.delphix_target.instance_id
}

output "delphix-tcw-tooling-oracle_system" {
  value = module.tooling.instance_id
}

output "delphix-tcw-virtualizationengine_system" {
  value = module.virtualization_engine.instance_id
}

output "delphix-tcw-maskingengine_system" {
  value = module.masking_engine.instance_id
}

output "delphix-tcw-jumpbox_system" {
  value = module.guacamole.instance_id
}

output "delphix-tcw-centos7-daf-app-prod_system" {
  value = module.prod_web_server.instance_id
}

output "delphix-tcw-centos7-daf-app-dev_system" {
  value = module.dev_web_server.instance_id
}

output "delphix-tcw-centos7-daf-app-test_system" {
  value = module.test_web_server.instance_id
}




