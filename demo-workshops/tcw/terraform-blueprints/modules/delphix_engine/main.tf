data "aws_ami" "de_ami" {
  most_recent = true
  filter {
    name = "name"
    values = ["${var.ami_name}"]
  }

  #From Delphix
  owners = ["${var.staged == "false" ? "180093685553" : "self"}"]
  // owners = ["self"]
}

resource "aws_instance" "delphix_engine" {
  instance_type = "t2.xlarge"
  # Lookup the correct AMI based on the region
  # we specified
  ami = "${data.aws_ami.de_ami.id}"
  count = "${var.how_many}"

  key_name = "${var.key_name}"
  associate_public_ip_address = "${var.associate_public_ip_address == "false" ? false : true}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${var.security_group_id}"]

  subnet_id = "${element(var.subnet_id, count.index)}"
  private_ip = "10.0.${count.index + 1}.${var.last_octet}"
  
  // ebs_optimized = true
  root_block_device {
      volume_type           = "gp2"
      volume_size           = "150"
      delete_on_termination = true
    }
  ebs_block_device {
        device_name = "/dev/sdb"
        volume_type = "gp2"
        volume_size = "8"
        delete_on_termination = true
    }
  ebs_block_device {
        device_name = "/dev/sdc"
        volume_type = "gp2"
        volume_size = "8"
        delete_on_termination = true
    }
  ebs_block_device {
        device_name = "/dev/sdd"
        volume_type = "gp2"
        volume_size = "8"
        delete_on_termination = true
    }
  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-DE-${count.index + 1}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

output "public_ip" {
  value = "${aws_instance.delphix_engine.*.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.delphix_engine.*.id}"
}

output "private_ip" {
  value = "${aws_instance.delphix_engine.*.private_ip}"
}
