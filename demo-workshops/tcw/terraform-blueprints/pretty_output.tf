output "environment" {
  value = formatlist(
    "\nJumpbox - Public IP: %s \n Access via browser @ http://%s:8080/labs\n",
    module.guacamole.public_ip,
    module.guacamole.public_ip,
  )
}

