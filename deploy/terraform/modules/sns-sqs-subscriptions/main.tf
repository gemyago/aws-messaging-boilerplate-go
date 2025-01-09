locals {
  sqs_subscriptions = [
    for index, subscription in var.subscriptions : {
      key               = subscription.subscription_id
      topic_arn         = subscription.topic_arn
      max_receive_count = subscription.max_receive_count
    }
  ]
}

resource "aws_sns_topic_subscription" "sqs_subscription" {
  for_each = {
    for index, subscription in local.sqs_subscriptions :
    subscription.key => subscription
  }
  topic_arn            = each.value.topic_arn
  protocol             = "sqs"
  endpoint             = aws_sqs_queue.primary[each.key].arn
  raw_message_delivery = true
}

resource "aws_sqs_queue" "primary" {
  for_each = {
    for index, subscription in local.sqs_subscriptions :
    subscription.key => subscription
  }
  name = "${var.resources_prefix}${var.app_name}-${each.value.key}"
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dead_letter[each.key].arn
    maxReceiveCount     = 3
  })
}

# allow SNS to send messages to the queue
resource "aws_sqs_queue_policy" "allow_sns_to_send_messages" {
  for_each = {
    for index, subscription in local.sqs_subscriptions :
    subscription.key => subscription
  }
  queue_url = aws_sqs_queue.primary[each.key].id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "sns.amazonaws.com"
        },
        Action   = "sqs:SendMessage",
        Resource = aws_sqs_queue.primary[each.key].arn,
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = each.value.topic_arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue" "dead_letter" {
  for_each = {
    for index, subscription in local.sqs_subscriptions :
    subscription.key => subscription
  }
  name = "${var.resources_prefix}${var.app_name}-${each.value.key}-dlq"
}