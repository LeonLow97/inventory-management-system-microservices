variable "name_prefix" {
  type = string
}

variable "vpc_cidr" {
  type = string
}

variable "availability_zones" {
  type = list(string)

  validation {
    condition     = length(var.availability_zones) == 2
    error_message = "The portfolio stack expects exactly two Availability Zones."
  }
}

variable "public_subnet_cidrs" {
  type = list(string)

  validation {
    condition     = length(var.public_subnet_cidrs) == 2
    error_message = "Provide one public subnet CIDR for each Availability Zone."
  }
}

variable "app_subnet_cidrs" {
  type = list(string)

  validation {
    condition     = length(var.app_subnet_cidrs) == 2
    error_message = "Provide one application subnet CIDR for each Availability Zone."
  }
}

variable "data_subnet_cidrs" {
  type = list(string)

  validation {
    condition     = length(var.data_subnet_cidrs) == 2
    error_message = "Provide one data subnet CIDR for each Availability Zone."
  }
}

variable "common_tags" {
  type = map(string)
}
