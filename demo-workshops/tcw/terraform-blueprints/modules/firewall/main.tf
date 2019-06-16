data "http" "your_ip" {
  url = "http://ipv4.icanhazip.com"

  # Optional request headers
  request_headers {
    "Accept" = "application/json"
  }
}

locals {
  cidr_addr = "${chomp("${data.http.your_ip.body}")}/32"
  
  default_firewall_ingress_cidr_blocks = [
    "${local.cidr_addr}"
  ]
}

resource "aws_security_group" "jumpbox" {
  name = "${var.name}-jumpbox"
  description = "Allow inbound traffic to jumpbox"
  vpc_id = "${var.vpc_id}"
  ingress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      #Remove the ', "0.0.0.0/0"' from the list below, to eliminate "Any In"
      cidr_blocks = ["${concat(local.default_firewall_ingress_cidr_blocks,var.addtl_firewall_ingress_cidr_blocks)}"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.name}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

resource "aws_security_group" "landshark" {
  name = "${var.name}"
  description = "Allow all inbound traffic"
  vpc_id = "${var.vpc_id}"

  ingress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      
      cidr_blocks = ["10.0.0.0/16", "0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.name}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

output "id" {
  value = "${aws_security_group.landshark.id}"
}

output "jumpbox" {
  value = "${aws_security_group.jumpbox.id}"
}