# In real world scenario the bus itself will be very likely 
# provisioned by a separate "infrastructure" like project. 
# Each "consumer" project will then just have "attachments" to it.
resource "aws_cloudwatch_event_bus" "event_bus" {
  name        = "${var.resources_prefix}app-events"
  description = "Example event bus to play around with. ${var.resources_description}"
}

module "event_bus_http_targets" {
  source                = "../modules/event-bus-http-targets"
  app_name              = var.app_name
  bus_name              = aws_cloudwatch_event_bus.event_bus.name
  aws_primary_region    = var.aws_primary_region
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description

  http_targets = [
    {
      event_source = "aws-sqs-boilerplate-go",
      detail_type  = "message",
      endpoint     = "https://webhook.site/6011e29c-7f67-4865-9c30-0c7c3b13fca4",
      method       = "POST",
      max_rps      = 20
    }
  ]
}