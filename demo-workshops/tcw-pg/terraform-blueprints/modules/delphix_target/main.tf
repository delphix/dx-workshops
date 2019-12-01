data "aws_ami" "delphix-ready-ami" {
  most_recent = true
  owners = ["self"]
  filter {
    name = "name"
    values = ["${var.ami_name}"]
  }
}

resource "aws_instance" "target" {
  count = "${var.how_many}"
  ami = "${data.aws_ami.delphix-ready-ami.id}"
  instance_type = "t2.small"
  key_name = "${var.key_name}"
  associate_public_ip_address = "${var.associate_public_ip_address == "false" ? false : true}"
  vpc_security_group_ids = "${var.security_group_id}"
  subnet_id = "${element(var.subnet_id, count.index)}"
  private_ip = "10.0.${count.index + 1}.${var.last_octet}"
  

  root_block_device {
    volume_type           = "gp2"
    volume_size           = "50"
    delete_on_termination = "true"
  }

  connection {
    type = "ssh"
    user = "centos"
    private_key = "${file("${var.key_name}.pem")}"
    timeout = "3m"
  }


  #Instance tags
  tags = "${merge(
    data.aws_ami.delphix-ready-ami.tags,
    var.default_tags,
    {
      "Name" = "${var.project}-db-${var.env_name}-${count.index + 1}",
      "STUDENT" = "${count.index + 1}",
      "source" = "${data.aws_ami.delphix-ready-ami.name}"
    }
  )}"
}

output "public_ip" {
  value = "${aws_instance.target.*.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.target.*.id}"
}

output "private_ip" {
  value = "${aws_instance.target.*.private_ip}"
}
