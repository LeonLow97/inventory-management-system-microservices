output "api_gateway_instance_ids" {
  value = {
    for availability_zone, instance in aws_instance.api_gateway : availability_zone => instance.id
  }
}

output "ecr_repository_url" {
  value = aws_ecr_repository.this.repository_url
}
