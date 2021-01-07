variable "aws_profile" {
  type        = string
  default     = "mtmo-non-prod"
  description = "The aws profile to use"
}

variable "mtmo_prod_aws_profile" {
  type        = string
  default     = "mtmo-prod"
  description = "The aws profile for MTMO"
}

variable "env_dns" {
  type        = string
  default     = "mtmostaging.com"
  description = "The environment base domain"
}
