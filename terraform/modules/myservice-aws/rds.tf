resource "aws_db_subnet_group" "default" {
  name       = "myservice-${var.env}"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "myservice-${var.env}"
  }
}

resource "random_password" "db_password" {
  length  = 16
  special = false
}

resource "aws_db_instance" "default" {
  identifier     = "myservice-${var.env}"
  engine         = "postgres"
  engine_version = "18"
  instance_class = var.rds_instance_class
  multi_az       = var.rds_multi_az

  db_name  = "myservice"
  username = "myservice"
  password = random_password.db_password.result

  db_subnet_group_name   = aws_db_subnet_group.default.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_encrypted     = true

  backup_retention_period = 7
  skip_final_snapshot     = false
  final_snapshot_identifier = "myservice-${var.env}-final"
  deletion_protection       = true

  performance_insights_enabled = true

  tags = {
    Name = "myservice-${var.env}"
  }
}
