module "vpc" {
  source = "../../modules/vpc"
  az1    = local.az1
  az2    = local.az2
}

module "security_group" {
  source = "../../modules/security_group"

  vpc_id            = module.vpc.vpc_id
  port_api_gateway  = var.port_api_gateway
  port_auth_service = var.port_auth_service
  port_postgresql   = var.port_postgresql
  public_ip         = "${local.public_ip}/32"
}

module "eks_security_group" {
  source = "../../modules/eks_security_group"

  vpc_id = module.vpc.vpc_id
}

module "rds" {
  source = "../../modules/rds"

  vpc_id              = module.vpc.vpc_id
  private_rds_subnets = module.vpc.private_rds_subnets
  db_sg               = module.security_group.db_sg_id
}

module "eks" {
    source = "../../modules/eks"

    vpc_id = module.vpc.vpc_id
    private_subnets = module.vpc.private_subnets
}

module "alb" {
  source = "../../modules/alb"

  alb_sg                      = module.security_group.alb_sg_id
  public_subnets              = module.vpc.public_subnets
  vpc_id                      = module.vpc.vpc_id
  api_gateway_instance_id_az1 = module.ec2.api_gateway_instance_id_az1
  api_gateway_instance_id_az2 = module.ec2.api_gateway_instance_id_az2
}
