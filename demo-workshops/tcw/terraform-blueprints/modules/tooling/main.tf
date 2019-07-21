data "aws_ami" "tooling_ami" {
  most_recent = true
  filter {
    name = "name"
    // values = ["delphix-CentOS7-tooling-*"]
    values = ["${var.ami_name}"]
  }
  owners = ["self"]
}

resource "aws_instance" "tooling" {
  instance_type = "t3.2xlarge"
  # Lookup the correct AMI based on the region
  # we specified
  count = "${var.how_many}"
  ami = "${data.aws_ami.tooling_ami.id}"
  connection {
    type = "ssh"
    user = "centos"
    private_key = "${file("${var.key_name}.pem")}"
    timeout = "3m"
  }

  associate_public_ip_address = true
  
  root_block_device {
    volume_type           = "gp2"
    volume_size           = "25"
    delete_on_termination = "true"
  }
  # The name of our SSH keypair you've created and downloaded
  # from the AWS console.
  #
  # https://console.aws.amazon.com/ec2/v2/home?region=us-west-2#KeyPairs:
  #
  key_name = "${var.key_name}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = "${var.security_group_id}"
  subnet_id = "${element(var.subnet_id, count.index)}"
  private_ip = "10.0.${count.index + 1}.${var.last_octet}"

 #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-TOOLING-${count.index + 1}",
      "STUDENT","${count.index + 1}"
      )
  )}"

}

output "public_ip" {
  value = "${aws_instance.tooling.*.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.tooling.*.id}"
}

output "private_ip" {
  value = "${aws_instance.tooling.*.private_ip}"
}