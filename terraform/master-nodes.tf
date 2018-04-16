resource "aws_instance" "master-node" {
  depends_on      = ["aws_nat_gateway.private-nat", "aws_route.private_route"]

  ami             = "${data.aws_ami.centos.id}"
  instance_type   = "${var.NodeTypes["Master"]}"
  key_name        = "${var.SshKey}"
  user_data       = "${file("scripts/init.sh")}"

  vpc_security_group_ids = ["${aws_security_group.master-sg.id}", "${aws_security_group.etcd-sg.id}", "${aws_security_group.allow-internal.id}"]
  subnet_id = "${aws_subnet.subnet-private-1.id}"

  count = "${var.Counts["Master"]}"

  root_block_device {
    volume_type = "gp2"
    volume_size = 50
  }

  lifecycle {
    create_before_destroy = true
  }

  tags {
    Type = "master"
    Name = "${var.ProjectName} - Master Node ${count.index + 1}"
    Project = "${var.ProjectName}"
    ProjectId = "${var.ProjectId}"
  }
}

resource "aws_lb_target_group_attachment" "master-to-master-lb1" {
  target_group_arn = "${aws_lb_target_group.master-lb-tg1.arn}"
  target_id        = "${aws_instance.master-node.*.id[count.index]}"
  port             = 8443

  count = "${var.Counts["Master"]}"
}

resource "aws_lb_target_group_attachment" "master-to-master-lb2" {
  target_group_arn = "${aws_lb_target_group.master-lb-tg2.arn}"
  target_id        = "${aws_instance.master-node.*.id[count.index]}"
  port             = 8443

  count = "${var.Counts["Master"]}"
}
