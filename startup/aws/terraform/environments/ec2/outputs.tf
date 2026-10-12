output "api_base_url" {
  description = "Public URL for the API Gateway."
  value       = "http://${module.load_balancer.dns_name}"
}

output "alb_dns_name" {
  description = "DNS name of the public Application Load Balancer."
  value       = module.load_balancer.dns_name
}

output "ecr_repository_url" {
  description = "ECR repository URL used by the image publishing helper."
  value       = module.applications.ecr_repository_url
}

output "aws_region" {
  description = "AWS region hosting the stack."
  value       = var.region
}

output "api_gateway_target_group_arn" {
  description = "Target group used by the deployment readiness check."
  value       = module.load_balancer.target_group_arn
}

output "availability_zones" {
  description = "Availability Zones selected for the stack."
  value       = local.availability_zones
}
