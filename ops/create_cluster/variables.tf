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
