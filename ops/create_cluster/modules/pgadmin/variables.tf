variable postgres_endpoint {
  type = string
}

variable "env_dns" {
  type        = string
  default     = "mtmostaging.com"
  description = "The environment base domain"
}
