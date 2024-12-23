resource "aws_sqs_queue" "messages" {
  name = "${var.resources_prefix}messages-queue"
}