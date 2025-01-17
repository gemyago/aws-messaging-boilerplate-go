terraform {
  backend "s3" {
    bucket     = "terraform-local"
    region     = "us-east-1"
    access_key = "local"
    secret_key = "local"
    endpoints = {
      s3  = "http://s3.localhost.localstack.cloud:4566"
      sts = "http://localhost:4566"
    }
  }
}
