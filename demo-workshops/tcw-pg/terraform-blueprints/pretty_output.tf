output "environment" {
  value = formatlist(
    "Jumpbox - Public IP: %s  Access via browser @ http://%s:8080/labs",
    module.jumpbox.public_ip,
    module.jumpbox.public_ip,
  )
}

