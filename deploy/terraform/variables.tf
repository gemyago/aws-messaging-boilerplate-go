variable "app_name" {
  type = string
}

variable "aws_primary_region" {
  type = string
}

variable "ci_mode" {
  type        = bool
  description = "Value is true if running in CI mode. Should be used in very specific cases."
}