# Event bus
resource "aws_cloudwatch_event_bus" "event_bus" {
  name = var.bus_name
}

resource "aws_cloudwatch_event_connection" "target_connection" {
  name               = "${var.bus_name}-target-connection"
  description        = "Connection to the target"
  authorization_type = "API_KEY"

  auth_parameters {
    api_key {
      key   = "Authorization"
      value = "Bearer NOT-USED"
    }
  }
}

resource "aws_cloudwatch_event_rule" "custom_source_events" {
  name        = "capture-custom-source-events"
  description = "Capture each AWS Console Sign In"
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name

  event_pattern = jsonencode({
    source: ["my.custom.source"],
    detail-type = [ "myDetailType" ]
  })
}

resource "aws_cloudwatch_event_api_destination" "test" {
  name                             = "api-destination"
  description                      = "An API Destination"
  invocation_endpoint              = "http://host.docker.internal:8080/messages/process"
  http_method                      = "POST"
  invocation_rate_limit_per_second = 20
  connection_arn                   = aws_cloudwatch_event_connection.target_connection.arn
}

resource "aws_cloudwatch_event_target" "test_target" {
  rule      = aws_cloudwatch_event_rule.custom_source_events.name
  arn       = aws_cloudwatch_event_api_destination.test.arn
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name
  # input_path = "$.detail"

  http_target {
    header_parameters = {
      "X-Message-ID" = "$.detail.id"
    }
  }
}