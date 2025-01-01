variable "resources_prefix" {
  type        = string
  description = "Prefix resources with given string. Useful to avoid name conflicts or setup test resources."
  default     = ""
}

variable "resources_description" {
  type        = string
  description = "Resources that support description field will have this value added."
  default     = ""
}

variable "local_stack_env" {
  type        = bool
  description = "Local stack environment"
}