variable "app_name" {
  type = string
}

variable "ci_mode" {
  type        = bool
  description = "Value is true if running in CI mode. Should be used in very specific cases."
}

variable "aws_primary_region" {
  type = string
}

variable "aws_access_key" {
  description = "AWS access key"
  type        = string
  default     = null
}

variable "aws_secret_key" {
  description = "AWS secret key"
  type        = string
  default     = null
}

variable "aws_endpoints" {
  description = "AWS service endpoints"
  type        = object({
    apigateway     = string
    apigatewayv2   = string
    cloudformation = string
    cloudwatch     = string
    dynamodb       = string
    ec2            = string
    es             = string
    elasticache    = string
    firehose       = string
    iam            = string
    kinesis        = string
    lambda         = string
    rds            = string
    redshift       = string
    route53        = string
    s3             = string
    secretsmanager = string
    ses            = string
    sns            = string
    sqs            = string
    ssm            = string
    stepfunctions  = string
    sts            = string
  })
  default     = null
}