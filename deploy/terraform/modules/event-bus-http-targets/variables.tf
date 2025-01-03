variable "app_name" {
  type = string
}

variable "bus_name" {
  type        = string
  description = "Name of the event bus"
}

variable "aws_primary_region" {
  type        = string
  description = "Primary region for the resources"
}

variable "resources_prefix" {
  type        = string
  description = "Prefix resources with given string. Useful to avoid name conflicts."
  default     = ""
}

variable "resources_description" {
  type        = string
  description = "Resources that support description field will have this value added."
  default     = ""
}

variable "http_targets" {
  type = list(object({
    target_id    = string
    event_source = string
    detail_type  = string
    endpoint     = string
    method       = string
    max_rps      = optional(number, 20)
  }))
  description = "List of HTTP targets to create"
}