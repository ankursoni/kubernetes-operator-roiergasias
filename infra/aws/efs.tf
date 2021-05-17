resource "aws_efs_file_system" "efs01" {
  creation_token = "${var.prefix}-${var.environment}-efs01"
  encrypted      = true

  tags = {
    Name = "${var.prefix}-${var.environment}-efs01"
  }
}

resource "aws_efs_mount_target" "emt01" {
  file_system_id  = aws_efs_file_system.efs01.id
  subnet_id       = module.vpc.private_subnets[0]
  security_groups = [aws_security_group.sg01.id]
}

resource "aws_efs_mount_target" "emt02" {
  file_system_id  = aws_efs_file_system.efs01.id
  subnet_id       = module.vpc.private_subnets[1]
  security_groups = [aws_security_group.sg01.id]
}

output "efs_id" {
  value = aws_efs_file_system.efs01.id
}
