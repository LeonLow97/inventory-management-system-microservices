variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "port_api_gateway" {
  type = number
}

variable "port_auth_service" {
  type = number
}

variable "port_postgresql" {
  type = number
}

variable "public_ip" {
  type = string
}
