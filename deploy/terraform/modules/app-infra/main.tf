resource "aws_cloudwatch_event_bus" "app_events" {
  name = "${var.resources_prefix}app-events"

  # The local stack has an issue of not assigning description. This leads
  # to a situation where plan is producing a constant diff on the description.
  # So we have to skip the description in local stack environment.
  # Maybe fixed by https://github.com/localstack/localstack/issues/12065
  description = var.local_stack_env ? "" : "Example event bus to play around with. ${var.resources_description}"
}

resource "aws_sns_topic" "dummy_messages" {
  name = "${var.resources_prefix}dummy-messages"
}