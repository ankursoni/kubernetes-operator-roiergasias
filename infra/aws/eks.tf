resource "aws_eks_cluster" "eks01" {
  name     = "${var.prefix}-${var.environment}-eks01"
  role_arn = aws_iam_role.iamr01.arn
  version  = "1.19"

  vpc_config {
    subnet_ids = [
      module.vpc.private_subnets[0], module.vpc.private_subnets[1],
      module.vpc.public_subnets[0], module.vpc.public_subnets[1]
    ]
  }

  depends_on = [
    aws_iam_role_policy_attachment.iamrpa01,
    aws_iam_role_policy_attachment.iamrpa02
  ]
}

resource "aws_iam_role" "iamr01" {
  name = "${var.prefix}-${var.environment}-iamr01"

  assume_role_policy = jsonencode({
    Version : "2012-10-17",
    Statement : [{
      Effect : "Allow",
      Principal : {
        Service : "eks.amazonaws.com"
      },
      Action : "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "iamrpa01" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.iamr01.name
}

resource "aws_iam_role_policy_attachment" "iamrpa02" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.iamr01.name
}


resource "aws_eks_node_group" "eksng01" {
  cluster_name    = aws_eks_cluster.eks01.name
  node_group_name = "${var.prefix}-${var.environment}-eksng01"
  node_role_arn   = aws_iam_role.iamr02.arn
  subnet_ids      = [for s in module.vpc.private_subnets : s]

  scaling_config {
    desired_size = var.node_count
    max_size     = var.node_count + 1
    min_size     = 1
  }

  launch_template {
    name    = aws_launch_template.lt01.name
    version = "$Latest"
  }

  lifecycle {
    ignore_changes = [
      scaling_config[0].desired_size,
      scaling_config[0].max_size,
      scaling_config[0].min_size,
    ]
  }

  depends_on = [
    aws_iam_role_policy_attachment.iamrpa03,
    aws_iam_role_policy_attachment.iamrpa04,
    aws_iam_role_policy_attachment.iamrpa05,
  ]
}

resource "aws_launch_template" "lt01" {
  name          = "${var.prefix}-${var.environment}-lt01"
  instance_type = "t2.medium"
}

resource "aws_iam_role" "iamr02" {
  name = "${var.prefix}-${var.environment}-iamr02"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "iamrpa03" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.iamr02.name
}

resource "aws_iam_role_policy_attachment" "iamrpa04" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.iamr02.name
}

resource "aws_iam_role_policy_attachment" "iamrpa05" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.iamr02.name
}


resource "aws_eks_fargate_profile" "eksfp01" {
  cluster_name           = aws_eks_cluster.eks01.name
  fargate_profile_name   = "${var.prefix}-${var.environment}-eksfp01"
  pod_execution_role_arn = aws_iam_role.iamr03.arn
  subnet_ids             = [for s in module.vpc.private_subnets : s]

  selector {
    namespace = "kube-system"
  }
  selector {
    namespace = "kubernetes-dashboard"
  }
  selector {
    namespace = "default"
  }
  selector {
    namespace = "roiergasias"
  }

  depends_on = [
    aws_iam_role_policy_attachment.iamrpa06
  ]
}

resource "aws_iam_role" "iamr03" {
  name = "${var.prefix}-${var.environment}-iamr03"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "eks-fargate-pods.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "iamrpa06" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"
  role       = aws_iam_role.iamr03.name
}
