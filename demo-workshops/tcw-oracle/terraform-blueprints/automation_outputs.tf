output "delphix-tcw-source_ip" {
  value = module.delphix_source.public_ip
}

output "delphix-tcw-source-oracle12_image" {
  value = module.delphix_source.instance_id
}

output "delphix-tcw-target_ip" {
  value = module.delphix_target.public_ip
}

output "delphix-tcw-target-oracle12_image" {
  value = module.delphix_target.instance_id
}

output "delphix-tcw-tooling_ip" {
  value = module.tooling.public_ip
}

output "delphix-tcw-tooling-oracle12_image" {
  value = module.tooling.instance_id
}

output "delphix-tcw-virtualizationengine_ip" {
  value = module.virtualization_engine.public_ip
}

output "delphix-tcw-virtualizationengine-oracle12_image" {
  value = module.virtualization_engine.instance_id
}

output "delphix-tcw-maskingengine_ip" {
  value = module.masking_engine.public_ip
}

output "delphix-tcw-maskingengine-oracle12_image" {
  value = module.masking_engine.instance_id
}

output "delphix-tcw-jumpbox_ip" {
  value = module.jumpbox.public_ip
}

output "delphix-tcw-jumpbox-oracle12_image" {
  value = module.jumpbox.instance_id
}

output "delphix-tcw-centos7-patients-app-oracle12_image" {
  value = module.prod_web_server.instance_id
}

output "delphix-tcw-source-oracle12_system" {
  value = module.delphix_source.instance_id
}

output "delphix-tcw-target-oracle12_system" {
  value = module.delphix_target.instance_id
}

output "delphix-tcw-tooling-oracle12_system" {
  value = module.tooling.instance_id
}

output "delphix-tcw-virtualizationengine_system" {
  value = module.virtualization_engine.instance_id
}

output "delphix-tcw-maskingengine_system" {
  value = module.masking_engine.instance_id
}

output "delphix-tcw-jumpbox_system" {
  value = module.jumpbox.instance_id
}

output "delphix-tcw-centos7-patients-app-prod_system" {
  value = module.prod_web_server.instance_id
}

output "delphix-tcw-centos7-patients-app-dev_system" {
  value = module.dev_web_server.instance_id
}

output "delphix-tcw-centos7-patients-app-test_system" {
  value = module.test_web_server.instance_id
}




