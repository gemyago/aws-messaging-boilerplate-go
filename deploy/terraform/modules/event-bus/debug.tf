resource "aws_cloudwatch_event_rule" "cloudwatch_debug_rule" {
  lifecycle {
    ignore_changes = [
      tags, tags_all
    ]
  }

  name           = "${var.resources_prefix}cloudwatch-debug-rule"
  description    = "Forward EventBridge internal events to CloudWatch log groups. ${var.resources_description}"
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name

  event_pattern = jsonencode({
    account : [
      "026264649083"
    ],
    # source : [
    #   "aws.events",
    #   "aws.cloudwatch"
    # ]
  })
}

resource "aws_cloudwatch_log_group" "debug_log_group" {
  lifecycle {
    ignore_changes = [
      tags, tags_all
    ]
  }
  name = "${local.bus_name}-debug-log-group"
}

resource "aws_cloudwatch_log_resource_policy" "debug_log_group_policy" {
  policy_name = "${local.bus_name}-debug-log-group-policy"
  policy_document = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Sid    = "TrustEventsToStoreLogEvent",
        Effect = "Allow",
        Principal = {
          Service = [
            "events.amazonaws.com",
            "delivery.logs.amazonaws.com"
          ]
        },
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "arn:aws:logs:${var.aws_primary_region}:${data.aws_caller_identity.current.account_id}:log-group:${aws_cloudwatch_log_group.debug_log_group.name}:*"
      }
    ]
  })
}

resource "aws_cloudwatch_event_target" "debug_target" {
  rule           = aws_cloudwatch_event_rule.cloudwatch_debug_rule.name
  arn            = aws_cloudwatch_log_group.debug_log_group.arn
  event_bus_name = aws_cloudwatch_event_bus.event_bus.name
}

