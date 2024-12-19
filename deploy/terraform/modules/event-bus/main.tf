# Event bus
resource "aws_cloudwatch_event_bus" "event_bus" {
  name = var.bus_name
}