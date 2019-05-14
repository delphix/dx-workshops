data "aws_ami" "guac_ami" {
  most_recent = true
  filter {
    name = "name"
    values = ["${var.ami_name}"]
  }
  owners = ["self"]
}

resource "aws_instance" "guacamole" {
  count = "${var.how_many}"
  ami = "${data.aws_ami.guac_ami.id}"
  instance_type = "t3.medium"
  key_name = "${var.key_name}"
  connection {
    type = "ssh"
    user = "ubuntu"
    private_key = "${file("${var.key_name}.pem")}"
    timeout = "3m"
  }
  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${var.security_group_id}"]
  subnet_id = "${element(var.subnet_id, count.index)}"
  private_ip = "10.0.${count.index + 1}.${var.last_octet}"
  associate_public_ip_address = true

  root_block_device {
    volume_size = 15
  }
  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-guacamole-${count.index + 1}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

output "public_ip" {
  value = "${aws_instance.guacamole.*.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.guacamole.*.id}"
}

output "private_ip" {
  value = "${aws_instance.guacamole.*.private_ip}"
}