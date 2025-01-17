provider "aws" {
  region = var.aws_primary_region

  # Set only for localstack. Otherwise use profile or AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
  # access_key = "local"
  # secret_key = "local"

  # Optionally set endpoints if running localstack
  # More details https://docs.localstack.cloud/user-guide/integrations/terraform/
  # endpoints {
  #   apigateway     = "http://localhost:4566"
  #   apigatewayv2   = "http://localhost:4566"
  #   cloudformation = "http://localhost:4566"
  #   cloudwatch     = "http://localhost:4566"
  #   eventbridge    = "http://localhost:4566"
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
}