variable "name_prefix" {
  type = string
}

variable "vpc_id" {
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

variable "common_tags" {
  type = map(string)
}
