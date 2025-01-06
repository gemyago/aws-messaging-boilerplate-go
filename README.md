# aws-sqs-boilerplate-go

[![Tests](https://github.com/gemyago/aws-sqs-boilerplate-go/actions/workflows/run-tests.yml/badge.svg)](https://github.com/gemyago/aws-sqs-boilerplate-go/actions/workflows/run-tests.yml)
[![Coverage](https://raw.githubusercontent.com/gemyago/aws-sqs-boilerplate-go/test-artifacts/coverage/golang-coverage.svg)](https://htmlpreview.github.io/?https://raw.githubusercontent.com/gemyago/aws-sqs-boilerplate-go/test-artifacts/coverage/golang-coverage.html)

Basic golang boilerplate for backend project that includes AWS SQS usage example.

Key features:
* [cobra](github.com/spf13/cobra) - CLI interactions
* [viper](github.com/spf13/viper) - Configuration management
* http.ServeMux is used as router (pluggable)
* uber [dig](go.uber.org/dig) is used as DI framework
  * for small projects it may make sense to setup dependencies manually
* `slog` is used for logs
* [slog-http](github.com/samber/slog-http) is used to produce access logs
* [testify](github.com/stretchr/testify) and [mockery](github.com/vektra/mockery) are used for tests
* [gow](github.com/mitranim/gow) is used to watch and restart tests or server

## Project structure

* [cmd/server](./cmd/server) is a main entrypoint to start API server
* [internal/api/http](./internal/api/http) - includes http routes related stuff
  * [internal/api/http/routes](./internal/api/http/routes) - add new routes here and register in [handler.go](./internal/api/http/server/handler.go)
* `internal/app` - is assumed to include application layer code (e.g business logic). Examples to be added.
* `internal/services` - lower level components are supposed to be here (e.g database access layer e.t.c). Examples to be added.

## Project Setup

Please have the following tools installed: 
* [direnv](https://github.com/direnv/direnv) 
* [pyenv](https://github.com/pyenv/pyenv?tab=readme-ov-file#installation)
* [gobrew](https://github.com/kevincobain2000/gobrew#install-or-update)

Python is required to run local setup script. 
```bash
# Install required python version
pyenv install -s

# Setup virtual environment
python -m venv .venv

# Reload direnv
direnv reload

# Install python dependencies
pip install -r requirements.txt
```

Install/Update go dependencies: 
```sh
# Install 
go mod download
make tools

# Update:
go get -u ./... && go mod tidy
```

### Setup LocalStack

LocalStack is used to run AWS services locally. To setup and provision the required resources, run the following command:

```bash
# Start LocalStack
docker compose up -d

# Provision resources
./scripts/localstack.py provision
make -C deploy/terraform init
make -C deploy/terraform plan
make -C deploy/terraform apply
```

### Lint and Tests

Run all lint and tests:
```bash
make lint
make test
```

Run specific tests:
```bash
# Run once
go test -v ./internal/api/http/routes/ --run TestHealthCheckRoutes

# Run same test multiple times
# This is useful for tests that are flaky
go test -v -count=5 ./internal/api/http/routes/ --run TestHealthCheckRoutes

# Run and watch
gow test -v ./internal/api/http/routes/ --run TestHealthCheckRoutes
```
### Run local API server:

```bash
# Regular mode
go run ./cmd/server/ start

# Watch mode (double ^C to stop)
gow run ./cmd/server/ start
```

## Deployment

This section describes how to deploy the application to AWS. Prior to deploying please make sure to initialize the AWS cli and configure the credentials. Please verify credentials by running the following command:
```bash
aws sts get-caller-identity
```

### Deployment

Deployment configuration is defined per environment and are stored in the [environments](./deploy/terraform/environments) directory. The `local` is a default environment that points to localstack and suitable for local development.

In order to create a new environment please create a new directory under the `environments` folder. Please name the directory according to the environment you are deploying to. If you do not wish to commit the configuration, please add `-local` suffix to the directory name (e.g `my-aws-local`). 

Use the [template](./deploy/terraform/environments/template) as a starting point for the new configuration. Update `backend.tf` and specify bucket name to store terraform state. Review and update other files as required, especially:
* `variables.tf`
  * Consider adding default values for `resources_prefix` and `resources_description`. This may be useful in a shared AWS account to distinguish resources.

For each new environment make sure the state bucket is available. You may use aws cli to create the bucket:
```bash
export bucket_name=<bucket_name>
export region=<region>
aws s3api create-bucket --bucket $bucket_name --region $region
aws s3api put-bucket-versioning --bucket $bucket_name --versioning-configuration Status=Enabled
unset bucket_name region
```
Make sure to pick globally unique bucket name. Example: `<aws-account>-<region>-terraform-state-<user>`

### Deploy

To deploy terraform configuration, run the following commands (from deploy/terraform directory):
```bash
# Set the environment to deploy to
# If not set, the default env is local which points
# to localstack
export DEPLOY_ENV=<env>

# Optionally cleanup terraform working directory.
# Obligatory if updating backend or changing aws credentials.
rm -r -f ${DEPLOY_ENV}/.terraform

make init

# Prepare and review the plan
make plan

# Make sure the plan looks good. Apply the plan
make apply

# Cleanup provisioned resources if needed
make plan_destroy
make apply

# Make sure to do it after the deployment
unset DEPLOY_ENV
```

## Useful commands

```bash
# Lock python dependencies (if updated)
pip freeze > requirements.txt

# Send custom event to AWS EventBridge
# Use awslocal to send event to localstack
aws events put-events --entries '[
  {
    "Source": "aws-sqs-boilerplate-go",
    "DetailType": "dummy-message",
    "Detail": "{\"id\": \"123\", \"name\": \"123\"}",
    "EventBusName": "app-events"
  }
]'
```

## Monitoring

Configured invocation rate and actual invocation rate. Alarm if exceeding, meaning scale-up is needed.