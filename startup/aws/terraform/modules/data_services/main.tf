resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-database"
  subnet_ids = var.data_subnet_ids

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-database"
  })
}

resource "aws_db_parameter_group" "this" {
  name   = "${var.name_prefix}-postgres"
  family = "postgres14"

  parameter {
    name         = "rds.force_ssl"
    value        = "0"
    apply_method = "immediate"
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-postgres"
  })
}

resource "aws_db_instance" "this" {
  identifier                 = "${var.name_prefix}-postgres"
  engine                     = "postgres"
  engine_version             = "14"
  instance_class             = "db.t3.micro"
  allocated_storage          = 20
  storage_type               = "gp3"
  storage_encrypted          = true
  db_subnet_group_name       = aws_db_subnet_group.this.name
  parameter_group_name       = aws_db_parameter_group.this.name
  vpc_security_group_ids     = [var.postgres_security_group_id]
  username                   = var.master_username
  password                   = var.master_password
  port                       = var.postgres_port
  publicly_accessible        = false
  multi_az                   = false
  backup_retention_period    = 0
  auto_minor_version_upgrade = true
  deletion_protection        = false
  skip_final_snapshot        = true
  copy_tags_to_snapshot      = true

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-postgres"
  })
}

resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.name_prefix}-redis"
  subnet_ids = var.data_subnet_ids

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-redis"
  })
}

resource "aws_elasticache_cluster" "this" {
  cluster_id               = "${var.name_prefix}-redis"
  engine                   = "redis"
  node_type                = "cache.t3.micro"
  num_cache_nodes          = 1
  port                     = var.redis_port
  subnet_group_name        = aws_elasticache_subnet_group.this.name
  security_group_ids       = [var.redis_security_group_id]
  snapshot_retention_limit = 0
  apply_immediately        = true

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-redis"
  })
}

resource "aws_s3_bucket" "schema" {
  bucket_prefix = "${var.name_prefix}-schema-"
  force_destroy = true

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-schema"
  })
}

resource "aws_s3_bucket_public_access_block" "schema" {
  bucket = aws_s3_bucket.schema.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "schema" {
  bucket = aws_s3_bucket.schema.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_object" "schema" {
  bucket       = aws_s3_bucket.schema.id
  key          = "users.sql"
  source       = var.schema_file
  etag         = filemd5(var.schema_file)
  content_type = "application/sql"
}

resource "aws_ssm_parameter" "database_name" {
  name  = "/authify/rds/db-name"
  type  = "String"
  value = var.database_name
}

resource "aws_ssm_parameter" "database_address" {
  name  = "/authify/rds/endpoint"
  type  = "String"
  value = aws_db_instance.this.address
}

resource "aws_ssm_parameter" "master_username" {
  name  = "/authify/rds/master-username"
  type  = "String"
  value = var.master_username
}

resource "aws_ssm_parameter" "authify_username" {
  name  = "/authify/rds/user"
  type  = "String"
  value = var.authify_username
}

resource "aws_ssm_parameter" "master_password" {
  name  = "/authify/rds/master-password"
  type  = "SecureString"
  value = var.master_password
}

resource "aws_ssm_parameter" "authify_password" {
  name  = "/authify/rds/authify-service-password"
  type  = "SecureString"
  value = var.authify_password
}

resource "aws_ssm_parameter" "redis_address" {
  name  = "/authify/elasticache/endpoint"
  type  = "String"
  value = aws_elasticache_cluster.this.cache_nodes[0].address
}

resource "aws_ssm_parameter" "schema_bucket" {
  name  = "/authify/s3/bucket/schema"
  type  = "String"
  value = aws_s3_bucket.schema.id
}

resource "aws_ssm_parameter" "jwt_secret" {
  name  = "/authify/auth/jwt-secret"
  type  = "SecureString"
  value = var.jwt_secret
}

resource "aws_ssm_parameter" "jwt_domain" {
  count = var.jwt_domain == "" ? 0 : 1

  name  = "/authify/auth/jwt-domain"
  type  = "String"
  value = var.jwt_domain
}
