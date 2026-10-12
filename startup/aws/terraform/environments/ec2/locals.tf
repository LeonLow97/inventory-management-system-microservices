data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

data "aws_ssm_parameter" "amazon_linux_2023" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
}

locals {
  availability_zones = slice(data.aws_availability_zones.available.names, 0, 2)

  app_subnets = {
    az1 = module.network.app_subnet_ids[0]
    az2 = module.network.app_subnet_ids[1]
  }

  common_tags = {
    Project     = "Authify"
    Environment = "portfolio"
    ManagedBy   = "Terraform"
  }
}
