variable "name_prefix" {
  type = string
}

variable "data_subnet_ids" {
  type = list(string)
}

variable "postgres_security_group_id" {
  type = string
}

variable "redis_security_group_id" {
  type = string
}

variable "postgres_port" {
  type = number
}

variable "redis_port" {
  type = number
}

variable "database_name" {
  type = string
}

variable "master_username" {
  type = string
}

variable "master_password" {
  type      = string
  sensitive = true
}

variable "authify_username" {
  type = string
}

variable "authify_password" {
  type      = string
  sensitive = true
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "jwt_domain" {
  type = string
}

variable "schema_file" {
  type = string
}

variable "common_tags" {
  type = map(string)
}
