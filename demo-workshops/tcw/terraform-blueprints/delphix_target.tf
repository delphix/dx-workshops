module "delphix_target" {
  source = "./modules/delphix_target"
  project = "${var.project}"
  key_name = "${var.key_name}"
  ami_name = "delphix-tcw-oracle12-target-${var.staged == "false" ? "unstaged" : "staged"}-*"
  how_many = "${var.how_many}"
  subnet_id = "${module.aws_subnet.id}"
  security_group_id = ["${module.aws_security_group.id}"]
  default_tags = "${local.default_tags}"
  env_name = "TARGET"
  last_octet = "30"
  associate_public_ip_address = "${var.associate_public_ip_address}"
}

output "Oracle Target" {
  value = "${
    formatlist(
      "\nOracle Target - Public IP: %s Private IP: %s\n",
      module.delphix_target.public_ip,
      module.delphix_target.private_ip
      )}"
}
