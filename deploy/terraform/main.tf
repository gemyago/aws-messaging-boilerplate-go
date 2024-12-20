provider "aws" {
  region = var.aws_primary_region

  access_key = var.aws_access_key
  secret_key = var.aws_secret_key

  dynamic "endpoints" {
    for_each = var.aws_endpoints != null ? [var.aws_endpoints] : []
    content {
      apigateway       = endpoints.value.apigateway
      apigatewayv2     = endpoints.value.apigatewayv2
      cloudformation   = endpoints.value.cloudformation
      cloudwatch       = endpoints.value.cloudwatch
      cloudwatchevents = endpoints.value.cloudwatch
      dynamodb         = endpoints.value.dynamodb
      ec2              = endpoints.value.ec2
      es               = endpoints.value.es
      elasticache      = endpoints.value.elasticache
      firehose         = endpoints.value.firehose
      iam              = endpoints.value.iam
      kinesis          = endpoints.value.kinesis
      lambda           = endpoints.value.lambda
      rds              = endpoints.value.rds
      redshift         = endpoints.value.redshift
      route53          = endpoints.value.route53
      s3               = endpoints.value.s3
      secretsmanager   = endpoints.value.secretsmanager
      ses              = endpoints.value.ses
      sns              = endpoints.value.sns
      sqs              = endpoints.value.sqs
      ssm              = endpoints.value.ssm
      stepfunctions    = endpoints.value.stepfunctions
      sts              = endpoints.value.sts
    }
  }
}

// add messages queue
resource "aws_sqs_queue" "messages" {
  name = "${var.resources_prefix}messages-queue"
}

// add event bus
module "event_bus" {
  source                = "./modules/event-bus"
  bus_name              = "messages-bus"
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description
}