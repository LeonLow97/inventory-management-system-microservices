output "alb_security_group_id" {
  value = aws_security_group.alb.id
}

output "api_gateway_security_group_id" {
  value = aws_security_group.api_gateway.id
}

output "authify_service_security_group_id" {
  value = aws_security_group.authify_service.id
}

output "postgres_security_group_id" {
  value = aws_security_group.postgres.id
}

output "redis_security_group_id" {
  value = aws_security_group.redis.id
}
