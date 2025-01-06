terraform {
  backend "s3" {
    bucket     = "<aws-account>-<region>-terraform-state-<user>"
    region     = "us-east-1"

    # Set only for localstack. Otherwise use profile or AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
    # access_key = ""
    # secret_key = ""
    
    # Optionally set endpoints if running localstack. You may need the below at least for the backend.
    # More details https://docs.localstack.cloud/user-guide/integrations/terraform/
    # endpoints = {
    #   s3  = "http://s3.localhost.localstack.cloud:4566"
    #   sts = "http://localhost:4566"
    # }
  }
}
