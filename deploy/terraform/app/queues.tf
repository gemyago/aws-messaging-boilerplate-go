# In real world scenario the topic may be provisioned by producer service
# In this case such topics should be specified as a variables
resource "aws_sns_topic" "messages" {
  name = "${var.resources_prefix}messages"
}

resource "aws_sns_topic_subscription" "messages" {
  topic_arn            = aws_sns_topic.messages.arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.messages.arn
  raw_message_delivery = true
}

resource "aws_sqs_queue" "messages" {
  name = "${var.resources_prefix}messages-queue"
}