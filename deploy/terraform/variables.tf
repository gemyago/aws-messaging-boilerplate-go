# Provided by caller makefile. Will be be set to repo name.
# See deploy/terraform/Makefile
variable "app_name" {
  type = string
}

# Common prefix for all resources. 
# Useful to avoid name conflicts in shared accounts.
# Please keep it short and alphanumeric with dash or underscore in the end.
# Example: "my-stuff-"
variable "resources_prefix" {
  type        = string
  description = "Prefix resources with given string. Useful to avoid name conflicts."
  default     = ""
}

# Description to add to resources that support it.
# Useful in shared accounts to identify the owner and the purpose.
# resources_description = "Provisioned by $USER"
variable "resources_description" {
  type        = string
  description = "Resources that support description field will have this value added."
  default     = ""
}

variable "ci_mode" {
  type        = bool
  description = "Value is true if running in CI mode. Should be used in very specific cases."
}

variable "aws_primary_region" {
  type        = string
  description = "Primary region for the resources"
  default     = "us-east-1"
}

# Either set the below variables or set the environment variables
# AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
# DO NOT commit real secrets
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

# If deploying to localstack, use the following values
# aws_endpoints = {
#   apigateway     = "http://localhost:4566"
#   apigatewayv2   = "http://localhost:4566"
#   cloudformation = "http://localhost:4566"
#   cloudwatch     = "http://localhost:4566"
#   dynamodb       = "http://localhost:4566"
#   ec2            = "http://localhost:4566"
#   es             = "http://localhost:4566"
#   elasticache    = "http://localhost:4566"
#   firehose       = "http://localhost:4566"
#   iam            = "http://localhost:4566"
#   kinesis        = "http://localhost:4566"
#   lambda         = "http://localhost:4566"
#   rds            = "http://localhost:4566"
#   redshift       = "http://localhost:4566"
#   route53        = "http://localhost:4566"
#   s3             = "http://s3.localhost.localstack.cloud:4566"
#   secretsmanager = "http://localhost:4566"
#   ses            = "http://localhost:4566"
#   sns            = "http://localhost:4566"
#   sqs            = "http://localhost:4566"
#   ssm            = "http://localhost:4566"
#   stepfunctions  = "http://localhost:4566"
#   sts            = "http://localhost:4566"
# }
variable "aws_endpoints" {
  description = "AWS service endpoints"
  type = object({
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
  default = null
}

variable "local_stack_env" {
  description = "Indicates if the environment is local stack. Used mainly to workaround local stack issues."
  type        = bool
  default     = false
}