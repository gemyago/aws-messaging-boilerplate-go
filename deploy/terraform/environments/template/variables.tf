# Provided by caller makefile. Will be be set to repo name.
# See deploy/terraform/Makefile
variable "app_name" {
  type = string
}

# Common prefix for all resources. 
# Useful to avoid name conflicts in shared accounts.
# Please keep it short and alphanumeric with dash or underscore in the end.
# Example: "my-stuff-"
variable "resources_prefix" {
  type        = string
  description = "Prefix resources with given string. Useful to avoid name conflicts."
  default     = ""
}

# Description to add to resources that support it.
# Useful in shared accounts to identify the owner and the purpose.
# resources_description = "Provisioned by $USER"
variable "resources_description" {
  type        = string
  description = "Resources that support description field will have this value added."
  default     = ""
}

variable "aws_primary_region" {
  type        = string
  description = "Primary region for the resources"
  default     = "us-east-1"
}