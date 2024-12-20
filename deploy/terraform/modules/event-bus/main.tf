locals {
  bus_name = "${var.resources_prefix}${var.bus_name}"
}

# Event bus
resource "aws_cloudwatch_event_bus" "event_bus" {
  name = local.bus_name
  description = "Event bus to route messages to the HTTP endpoint. ${var.resources_description}"
}

resource "aws_cloudwatch_event_connection" "target_connection" {
  name               = "${local.bus_name}-target-connection"
  description        = "Connection to the target. ${var.resources_description}"
  authorization_type = "API_KEY"

  auth_parameters {
    api_key {
      key   = "Authorization"
      value = "Bearer NOT-USED"
    }
  }
}

resource "aws_cloudwatch_event_rule" "custom_source_events" {
  name           = "${var.resources_prefix}capture-custom-source-events"
  description    = "Capture events from a custom source. ${var.resources_description}"
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name

  event_pattern = jsonencode({
    source : ["my.custom.source"],
    detail-type = ["myDetailType"]
  })
}

resource "aws_cloudwatch_event_api_destination" "test" {
  name                             = "${var.resources_prefix}api-destination"
  description                      = "An API Destination. ${var.resources_description}"
  invocation_endpoint              = "https://webhook.site/6011e29c-7f67-4865-9c30-0c7c3b13fca4"
  http_method                      = "POST"
  invocation_rate_limit_per_second = 20
  connection_arn                   = aws_cloudwatch_event_connection.target_connection.arn
}

resource "aws_cloudwatch_event_target" "test_target" {
  rule           = aws_cloudwatch_event_rule.custom_source_events.name
  arn            = aws_cloudwatch_event_api_destination.test.arn
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name
  # input_path = "$.detail"

  http_target {
    header_parameters = {
      "X-Message-ID" = "$.detail.id"
    }
  }
}