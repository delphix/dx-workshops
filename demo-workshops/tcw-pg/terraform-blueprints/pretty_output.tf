output "environment" {
  value = formatlist(
    "Jumpbox - Public IP: %s  Access via browser @ http://%s:8080/labs",
    module.guacamole.public_ip,
    module.guacamole.public_ip,
  )
}

