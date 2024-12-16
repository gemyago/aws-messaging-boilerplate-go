provider "aws" {
  access_key = "test"
  secret_key = "test"
  region     = "us-east-1"

  endpoints {
    s3 = "http://s3.localhost.localstack.cloud:4566"
  }
}