variable "name_prefix" {
  type = string
}

variable "region" {
  type = string
}

variable "ami_id" {
  type = string
}

variable "image_tag" {
  type = string
}

variable "ecr_repository_name" {
  type = string
}

variable "app_subnets" {
  type = map(string)
}

variable "api_gateway_security_group_id" {
  type = string
}

variable "authify_service_security_group_id" {
  type = string
}

variable "api_gateway_instance_profile" {
  type = string
}

variable "authify_service_instance_profile" {
  type = string
}

variable "api_gateway_port" {
  type = number
}

variable "authify_service_port" {
  type = number
}

variable "postgres_port" {
  type = number
}

variable "redis_port" {
  type = number
}

variable "api_gateway_instance_type" {
  type = string
}

variable "authify_service_instance_type" {
  type = string
}

variable "use_spot_instances" {
  type = bool
}

variable "parameter_names" {
  type = object({
    database_name    = string
    database_address = string
    master_username  = string
    master_password  = string
    authify_username = string
    authify_password = string
    redis_address    = string
    jwt_secret       = string
    jwt_domain       = string
  })
}

variable "schema_bucket_name" {
  type = string
}

variable "schema_object_key" {
  type = string
}

variable "common_tags" {
  type = map(string)
}
