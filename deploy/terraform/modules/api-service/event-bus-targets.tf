module "event_bus_http_targets" {
  source                = "../event-bus-http-targets"
  app_name              = var.app_name
  bus_name              = var.app_bus_name
  aws_primary_region    = var.aws_primary_region
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description

  http_targets = [
    {
      target_id    = "process-dummy-messages",
      event_source = "aws-messaging-boilerplate-go",
      detail_type  = "dummy-message",
      endpoint     = "${var.service_endpoint}/dummy-messages/process",
      method       = "POST",
      max_rps      = 20
    }
  ]
}