# Event bus name
variable "bus_name" {
  type        = string
  description = "Name of the event bus"
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