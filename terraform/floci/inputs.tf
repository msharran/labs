variable "aws_endpoint_url" {
  description = "Floci AWS-compatible endpoint. Set by .envrc."
  type        = string
  default     = "http://localhost:4566"
}

variable "aws_region" {
  description = "Local AWS region. Set by .envrc."
  type        = string
  default     = "ap-south-1"
}

variable "network_cidr" {
  type        = string
  description = "network cidr block"
}
