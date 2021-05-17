module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.prefix}-${var.environment}-vpc01"
  cidr = "10.0.0.0/16"

  azs             = ["${var.region}a", "${var.region}b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = true
  enable_vpn_gateway = true

  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    "kubernetes.io/cluster/${var.prefix}-${var.environment}-eks01" = "shared"
  }
}

resource "aws_security_group" "sg01" {
  name   = "${var.prefix}-${var.environment}-sg01"
  vpc_id = module.vpc.vpc_id

  ingress {
    description = "EFS from VPC"
    from_port   = 2049
    to_port     = 2049
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}