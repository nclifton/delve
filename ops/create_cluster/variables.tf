variable "aws_profile" {
  type        = string
  default     = "mtmo-non-prod"
  description = "The aws profile to use"
}

variable "aws_region" {
  type        = string
  default     = "ap-southeast-2"
  description = "The aws region deploying to"
}

variable "mtmo_prod_aws_profile" {
  type        = string
  default     = "mtmo-prod"
  description = "The aws profile for MTMO"
}

variable "env_dns" {
  type        = string
  description = "The environment base domain"
}
