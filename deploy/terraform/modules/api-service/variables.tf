variable "app_name" {
  type = string
}

variable "app_bus_name" {
  type        = string
  description = "Name of the app event bus"
}

variable "dummy_messages_topic_arn" {
  type        = string
  description = "The ARN of the dummy messages SNS topic"
}

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

variable "aws_primary_region" {
  type = string
}

variable "service_endpoint" {
  type        = string
  description = "Base URL of the service"
}