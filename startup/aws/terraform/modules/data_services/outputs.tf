output "schema_bucket_arn" {
  value = aws_s3_bucket.schema.arn
}

output "schema_bucket_name" {
  value = aws_s3_bucket.schema.id
}

output "schema_object_key" {
  value = aws_s3_object.schema.key
}

output "parameter_names" {
  value = {
    database_name    = aws_ssm_parameter.database_name.name
    database_address = aws_ssm_parameter.database_address.name
    master_username  = aws_ssm_parameter.master_username.name
    master_password  = aws_ssm_parameter.master_password.name
    authify_username = aws_ssm_parameter.authify_username.name
    authify_password = aws_ssm_parameter.authify_password.name
    redis_address    = aws_ssm_parameter.redis_address.name
    jwt_secret       = aws_ssm_parameter.jwt_secret.name
    jwt_domain       = var.jwt_domain == "" ? "" : aws_ssm_parameter.jwt_domain[0].name
  }
}
