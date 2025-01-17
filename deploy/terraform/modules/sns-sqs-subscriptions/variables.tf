variable "app_name" {
  type = string
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

variable "subscriptions" {
  type = list(object({
    subscription_id   = string
    topic_arn         = string
    max_receive_count = optional(number, 3)
  }))
  description = "List of SNS subscriptions to SQS queues."
}