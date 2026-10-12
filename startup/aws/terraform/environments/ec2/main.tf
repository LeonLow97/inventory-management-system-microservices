module "network" {
  source = "../../modules/network"

  name_prefix         = var.name_prefix
  vpc_cidr            = var.vpc_cidr
  availability_zones  = local.availability_zones
  public_subnet_cidrs = var.public_subnet_cidrs
  app_subnet_cidrs    = var.app_subnet_cidrs
  data_subnet_cidrs   = var.data_subnet_cidrs
  common_tags         = local.common_tags
}

module "security_groups" {
  source = "../../modules/security_groups"

  name_prefix          = var.name_prefix
  vpc_id               = module.network.vpc_id
  api_gateway_port     = var.api_gateway_port
  authify_service_port = var.authify_service_port
  postgres_port        = var.postgres_port
  redis_port           = var.redis_port
  common_tags          = local.common_tags
}

module "data_services" {
  source = "../../modules/data_services"

  name_prefix                = var.name_prefix
  data_subnet_ids            = module.network.data_subnet_ids
  postgres_security_group_id = module.security_groups.postgres_security_group_id
  redis_security_group_id    = module.security_groups.redis_security_group_id
  postgres_port              = var.postgres_port
  redis_port                 = var.redis_port
  database_name              = var.database_name
  master_username            = var.master_username
  master_password            = var.master_password
  authify_username           = var.authify_username
  authify_password           = var.authify_password
  jwt_secret                 = var.jwt_secret
  jwt_domain                 = var.jwt_domain
  schema_file                = abspath(var.schema_file)
  common_tags                = local.common_tags
}

module "iam" {
  source = "../../modules/iam"

  name_prefix       = var.name_prefix
  region            = var.region
  account_id        = data.aws_caller_identity.current.account_id
  schema_bucket_arn = module.data_services.schema_bucket_arn
  common_tags       = local.common_tags
}

module "applications" {
  source = "../../modules/applications"

  name_prefix                       = var.name_prefix
  region                            = var.region
  ami_id                            = data.aws_ssm_parameter.amazon_linux_2023.value
  image_tag                         = var.image_tag
  ecr_repository_name               = var.ecr_repository_name
  app_subnets                       = local.app_subnets
  api_gateway_security_group_id     = module.security_groups.api_gateway_security_group_id
  authify_service_security_group_id = module.security_groups.authify_service_security_group_id
  api_gateway_instance_profile      = module.iam.api_gateway_instance_profile
  authify_service_instance_profile  = module.iam.authify_service_instance_profile
  api_gateway_port                  = var.api_gateway_port
  authify_service_port              = var.authify_service_port
  postgres_port                     = var.postgres_port
  redis_port                        = var.redis_port
  api_gateway_instance_type         = var.api_gateway_instance_type
  authify_service_instance_type     = var.authify_service_instance_type
  use_spot_instances                = var.use_spot_instances
  parameter_names                   = module.data_services.parameter_names
  schema_bucket_name                = module.data_services.schema_bucket_name
  schema_object_key                 = module.data_services.schema_object_key
  common_tags                       = local.common_tags

  depends_on = [module.iam]
}

module "load_balancer" {
  source = "../../modules/load_balancer"

  name_prefix              = var.name_prefix
  vpc_id                   = module.network.vpc_id
  public_subnet_ids        = module.network.public_subnet_ids
  security_group_id        = module.security_groups.alb_security_group_id
  api_gateway_instance_ids = module.applications.api_gateway_instance_ids
  api_gateway_port         = var.api_gateway_port
  common_tags              = local.common_tags
}
