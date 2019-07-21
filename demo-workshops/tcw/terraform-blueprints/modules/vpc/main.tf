resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
    
  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-vpc}"
      )
  )}"
}

resource "aws_route" "r"{
  route_table_id = "${aws_vpc.main.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.main.id}"
}

resource "aws_internet_gateway" "main" {
    vpc_id = "${aws_vpc.main.id}"

  #Instance tags
  tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-ig"
      )
  )}"
}

output "id" {
  value = "${aws_vpc.main.id}"
}

output "main_route_table_id" {
  value = "${aws_vpc.main.main_route_table_id}"
}