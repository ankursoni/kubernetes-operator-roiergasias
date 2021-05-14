resource "aws_s3_bucket" "s3b01" {
  bucket = "${var.prefix}-${var.environment}-s3b01"
  acl    = "private"

  tags = {
    Name = "${var.prefix}-${var.environment}-s3b01"
  }
}
