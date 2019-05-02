module "tooling" {
  source = "./modules/tooling"
  project = "${var.project}"
  key_name = "${var.key_name}"
  how_many = "${var.how_many}"
  subnet_id = "${module.aws_subnet.id}"
  security_group_id = ["${module.aws_security_group.id}"]
  default_tags = "${local.default_tags}"
  last_octet = "6"
  ami_name = "delphix-tcw-tooling-oracle-${var.staged == "false" ? "unstaged" : "staged"}-*"
}

output "TOOLING" {
  value = "${
    formatlist(
      "\nTooling - Public IP: %s Private IP: %s\n",
      module.tooling.public_ip,
      module.tooling.private_ip
      )}"
}