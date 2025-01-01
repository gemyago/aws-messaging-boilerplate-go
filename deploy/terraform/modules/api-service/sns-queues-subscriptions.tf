resource "aws_sns_topic_subscription" "dummy_messages" {
  topic_arn            = var.dummy_messages_topic_arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.messages.arn
  raw_message_delivery = true
}

resource "aws_sqs_queue" "dummy_messages" {
  name = "${var.resources_prefix}${var.app_name}-dummy-messages"
}