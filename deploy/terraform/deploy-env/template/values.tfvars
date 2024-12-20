# Adjust the below values as needed

# Primary region to deploy to
aws_primary_region = "us-east-1"

# Either set the below (DO NOT COMMIT) or set the environment variables
# AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
# aws_access_key     = "XXX"
# aws_secret_key     = "YYY"

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