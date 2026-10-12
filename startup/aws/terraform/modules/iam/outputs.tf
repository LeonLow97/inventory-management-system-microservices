output "api_gateway_instance_profile" {
  value = aws_iam_instance_profile.api_gateway.name
}

output "authify_service_instance_profile" {
  value = aws_iam_instance_profile.authify_service.name
}
