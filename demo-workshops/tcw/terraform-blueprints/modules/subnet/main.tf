resource "aws_subnet" "aw_sub" {
    count = "${var.how_many}"
    vpc_id = "${var.vpc_id}"
    availability_zone = "${var.availability_zone}"
    cidr_block = "${lookup(var.cidr_blocks, count.index + 1)}"
    map_public_ip_on_launch = true

    #Instance tags
    tags = "${merge(
    var.default_tags,
    map(
      "Name", "${var.project}-subnet-${count.index + 1}",
      "STUDENT","${count.index + 1}"
      )
  )}"
}

output "id" {
  value = "${aws_subnet.aw_sub.*.id}"
}