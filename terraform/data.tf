
data "aws_availability_zones" "frankfurt" {}


data "aws_ami" "centos" {
  most_recent = true


  filter {
    name   = "name"
    values = ["CentOS Linux 7 x86_64 HVM EBS *"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["679593333241"] # CentOS official
}

// AmazonElasticFileSystemReadOnlyAccess policy used for EFS
data "aws_iam_policy" "efs-read-policy" {
  arn = "arn:aws:iam::aws:policy/AmazonElasticFileSystemReadOnlyAccess"
}