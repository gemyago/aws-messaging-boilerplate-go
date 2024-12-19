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