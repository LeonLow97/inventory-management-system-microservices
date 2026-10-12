variable "region" {
  description = "AWS region for the Authify stack."
  type        = string
  default     = "ap-southeast-1"
}

variable "name_prefix" {
  description = "Prefix for AWS resource names."
  type        = string
  default     = "authify"

  validation {
    condition     = can(regex("^[a-z][a-z0-9-]{2,19}$", var.name_prefix))
    error_message = "name_prefix must be 3 to 20 lowercase letters, numbers, or hyphens, and start with a letter."
  }
}

variable "vpc_cidr" {
  description = "CIDR range for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDRs, one for each selected Availability Zone."
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "app_subnet_cidrs" {
  description = "Private application subnet CIDRs, one for each selected Availability Zone."
  type        = list(string)
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "data_subnet_cidrs" {
  description = "Private database and cache subnet CIDRs, one for each selected Availability Zone."
  type        = list(string)
  default     = ["10.0.5.0/24", "10.0.6.0/24"]
}

variable "database_name" {
  description = "Database created for the Authify Service."
  type        = string
  default     = "authify_users_db"
}

variable "schema_file" {
  description = "Path to the Authify Service database schema. The Make targets supply an absolute path."
  type        = string

  validation {
    condition     = fileexists(var.schema_file)
    error_message = "schema_file must point to an existing SQL file."
  }
}

variable "master_username" {
  description = "RDS administrator username used during database initialization."
  type        = string
  default     = "authifymaster"
}

variable "master_password" {
  description = "RDS administrator password. Keep terraform.tfvars private."
  type        = string
  sensitive   = true

  validation {
    condition     = length(var.master_password) >= 12
    error_message = "master_password must contain at least 12 characters."
  }
}

variable "authify_username" {
  description = "Least-privilege PostgreSQL role used by Authify Service."
  type        = string
  default     = "authifyuser"
}

variable "authify_password" {
  description = "Password for the Authify Service PostgreSQL role. Keep terraform.tfvars private."
  type        = string
  sensitive   = true

  validation {
    condition     = length(var.authify_password) >= 12
    error_message = "authify_password must contain at least 12 characters."
  }
}

variable "jwt_secret" {
  description = "Signing secret shared by API Gateway and Authify Service. Keep terraform.tfvars private."
  type        = string
  sensitive   = true

  validation {
    condition     = length(var.jwt_secret) >= 32
    error_message = "jwt_secret must contain at least 32 characters."
  }
}

variable "jwt_domain" {
  description = "Optional cookie domain. Leave empty to use a host-only cookie."
  type        = string
  default     = ""
}

variable "ecr_repository_name" {
  description = "ECR repository shared by the API Gateway and Authify Service images."
  type        = string
  default     = "authify"
}

variable "image_tag" {
  description = "Immutable image release tag. The Make target supplies a unique default."
  type        = string
  default     = "portfolio"
}

variable "api_gateway_instance_type" {
  description = "EC2 instance type for API Gateway."
  type        = string
  default     = "t3.micro"
}

variable "authify_service_instance_type" {
  description = "EC2 instance type for Authify Service."
  type        = string
  default     = "t3.micro"
}

variable "use_spot_instances" {
  description = "Use Spot instances to reduce portfolio project costs. They can be interrupted."
  type        = bool
  default     = true
}

variable "api_gateway_port" {
  description = "HTTP port exposed by API Gateway."
  type        = number
  default     = 8080
}

variable "authify_service_port" {
  description = "gRPC port exposed by Authify Service."
  type        = number
  default     = 50051
}

variable "postgres_port" {
  description = "PostgreSQL port."
  type        = number
  default     = 5432
}

variable "redis_port" {
  description = "Redis port."
  type        = number
  default     = 6379
}
