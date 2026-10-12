variable "name_prefix" {
  type = string
}

variable "region" {
  type = string
}

variable "account_id" {
  type = string
}

variable "schema_bucket_arn" {
  type = string
}

variable "common_tags" {
  type = map(string)
}
