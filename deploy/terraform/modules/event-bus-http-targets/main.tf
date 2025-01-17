locals {
  http_targets = [
    for index, target in var.http_targets : {
      key          = target.target_id
      event_source = target.event_source
      detail_type  = target.detail_type
      endpoint     = target.endpoint
      method       = target.method
      max_rps      = target.max_rps
    }
  ]
}

data "aws_caller_identity" "current" {}

# TODO: This should transition to OAUTH2
# Private connection will be possible once below is done
# https://github.com/hashicorp/terraform-provider-aws/issues/40384
resource "aws_cloudwatch_event_connection" "target_connection" {
  name               = "${var.bus_name}-${var.app_name}-target-connection"
  description        = "Connection to the target. ${var.resources_description}"
  authorization_type = "API_KEY"

  auth_parameters {
    api_key {
      key   = "Authorization"
      value = "Bearer NOT-USED"
    }
  }
}

resource "aws_cloudwatch_event_rule" "capture_source_events" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }

  lifecycle {
    ignore_changes = [
      tags, tags_all
    ]
  }

  name_prefix    = "${var.resources_prefix}${var.app_name}-"
  description    = "Capture ${each.value.event_source} events of ${each.value.detail_type} type. ${var.resources_description}"
  event_bus_name = var.bus_name

  event_pattern = jsonencode({
    source : [each.value.event_source],
    detail-type = [each.value.detail_type]
  })
}

resource "random_id" "destination_name" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }

  byte_length = 4
  prefix      = "${var.resources_prefix}${var.app_name}-"

  keepers = {
    index = each.key
  }
}

resource "aws_cloudwatch_event_api_destination" "http_destination" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }
  name                             = random_id.destination_name[each.key].hex
  description                      = "An HTTP destination for ${each.value.event_source} events of ${each.value.detail_type} type. ${var.resources_description}"
  invocation_endpoint              = each.value.endpoint
  http_method                      = each.value.method
  invocation_rate_limit_per_second = each.value.max_rps
  connection_arn                   = aws_cloudwatch_event_connection.target_connection.arn
}

resource "aws_iam_role" "http_target_role" {
  lifecycle {
    ignore_changes = [
      tags, tags_all
    ]
  }

  description = "Role for the EventBridge ${var.bus_name} HTTP bus target. ${var.resources_description}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "events.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "http_target_role_policy" {
  role = aws_iam_role.http_target_role.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = "events:InvokeApiDestination",
        Resource = [
          for destination in aws_cloudwatch_event_api_destination.http_destination : destination.arn
        ]
      }
    ]
  })
}

resource "aws_cloudwatch_event_target" "api_target" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }
  rule           = aws_cloudwatch_event_rule.capture_source_events[each.key].name
  arn            = aws_cloudwatch_event_api_destination.http_destination[each.key].arn
  event_bus_name = var.bus_name
  role_arn       = aws_iam_role.http_target_role.arn
  input_path     = "$.detail"

  http_target {
    header_parameters = {
      # Below may not get substituted correctly in localstack setup
      # May be fixed with https://github.com/localstack/localstack/issues/12062
      "X-Message-ID"     = "$.id"
      "X-Message-Type"   = "$.detail-type"
      "X-Message-Source" = "$.source"
      "X-Message-Time"   = "$.time"
    }
  }

  dead_letter_config {
    arn = aws_sqs_queue.dead_letter[each.key].arn
  }

  retry_policy {
    maximum_retry_attempts       = 1 # TODO: This needs to be parameterized
    maximum_event_age_in_seconds = 180
  }
}

# allow event bridge rule to send messages to the dead letter queue
resource "aws_sqs_queue_policy" "allow_event_bridge_dlq" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }
  queue_url = aws_sqs_queue.dead_letter[each.key].id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "events.amazonaws.com"
        },
        Action   = "sqs:SendMessage",
        Resource = aws_sqs_queue.dead_letter[each.key].arn,
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = aws_cloudwatch_event_rule.capture_source_events[each.key].arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue" "dead_letter" {
  for_each = {
    for index, target in local.http_targets :
    target.key => target
  }
  name = "${var.resources_prefix}${var.app_name}-${each.value.key}-dlq"
}