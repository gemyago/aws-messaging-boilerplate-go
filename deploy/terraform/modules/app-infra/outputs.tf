output "dummy_messages_topic_arn" {
  value       = aws_sns_topic.dummy_messages.arn
  description = "The ARN of the dummy messages SNS topic"
}

output "app_bus_name" {
  value       = aws_cloudwatch_event_bus.app_events.name
  description = "The name of the app event bus"
}