variable "name_prefix" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "public_subnet_ids" {
  type = list(string)
}

variable "security_group_id" {
  type = string
}

variable "api_gateway_instance_ids" {
  type = map(string)
}

variable "api_gateway_port" {
  type = number
}

variable "common_tags" {
  type = map(string)
}
