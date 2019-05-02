data "aws_ami" "delphix-ready-ami" {
  most_recent = true
  owners = ["self"]
  filter {
    name = "name"
    values = ["${var.ami_name}"]
  }
}

resource "aws_instance" "web_server" {
  count = "${var.how_many}"
  ami = "${data.aws_ami.delphix-ready-ami.id}"
  instance_type = "t2.micro"
  key_name = "${var.key_name}"

  vpc_security_group_ids = ["${var.security_group_id}"]
  subnet_id = "${element(var.subnet_id, count.index)}"
  private_ip = "10.0.${count.index + 1}.${var.last_octet}"
  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-app-${var.env_name}-${count.index + 1}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

// output "student" {
//   value = "${aws_instance.web_server.*.tags.STUDENT ? aws_instance.web_server.*.tags.STUDENT : 0}"
// }
output "public_ip" {
  value = "${aws_instance.web_server.*.public_ip}"
}

output "instance_id" {
  value = "${aws_instance.web_server.*.id}"
}

output "private_ip" {
  value = "${aws_instance.web_server.*.private_ip}"
}