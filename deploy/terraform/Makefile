.PHONY: init check-deploy-env plan plan_destroy apply

DEPLOY_ENV ?= local

# taken from https://gist.github.com/sighingnow/deee806603ec9274fd47
# below may need to be extended to handle more cases (e.g M1)
UNAME_S := $(shell uname -s)
UNAME_P = $(shell uname -p)
ifeq ($(UNAME_S),Linux)
	os=linux
else
	os=darwin
endif
ifeq ($(UNAME_P),arm)
	arch=arm64
else
	arch=amd64
endif

tar_version := $(shell tar --version)
ifeq ($(findstring GNU tar,$(tar_version)),GNU tar)
	tar_ignore_arg=--ignore-failed-read
	tar_absolute_names=--absolute-names
else
	tar_ignore_arg=
	tar_absolute_names=
endif

bin=bin
tmp=./tmp
app_name:=$(shell basename `git rev-parse --show-toplevel`)

terraform_version=$(shell cat .terraform-version)
terraform_zip_name=terraform_$(terraform_version)_$(os)_$(arch).zip
terraform_zip=$(tmp)/$(terraform_zip_name)
terraform_val_file=environments/$(DEPLOY_ENV)/values.tfvars
terraform_backend_bucket_file=environments/$(DEPLOY_ENV)/state-bucket.s3.tfbackend
terraform_backend_prefix_file=environments/$(DEPLOY_ENV)/state-prefix.s3.tfbackend
terraform_vars=-var "app_name=$(app_name)"
terraform_plan_args=-input=false

ifdef CI
terraform=$(shell which terraform)
terraform_plan_args+=-no-color
else
terraform=$(bin)/terraform
endif


$(tmp):
	mkdir -p $(tmp)

$(bin):
	mkdir -p $(bin)

$(terraform_zip): .terraform-version $(tmp)
	curl https://releases.hashicorp.com/terraform/$(terraform_version)/$(terraform_zip_name) -o $(terraform_zip)

ifndef CI
$(terraform): $(terraform_zip) $(bin)
	unzip -o $(terraform_zip) -d $(bin)
	touch $(terraform)
endif

check-deploy-env:
ifndef DEPLOY_ENV
	$(error DEPLOY_ENV is undefined)
endif
ifeq ("$(wildcard $(terraform_val_file))","")
	$(error $(terraform_val_file) not exists)
endif

$(terraform_backend_prefix_file):
	echo 'key  = "$(app_name)-state"' > $@

init: $(terraform) check-deploy-env $(terraform_backend_prefix_file)
	$(terraform) init -backend-config=$(terraform_backend_bucket_file) -backend-config=$(terraform_backend_prefix_file)

init_upgrade: $(terraform) check-deploy-env $(terraform_backend_prefix_file)
	$(terraform) init -backend-config=$(terraform_backend_bucket_file) -backend-config=$(terraform_backend_prefix_file) -upgrade

import: $(terraform) check-deploy-env
	@echo $(terraform) import $(terraform_vars)

state: $(terraform) check-deploy-env
	@echo $(terraform) state

plan: $(terraform) check-deploy-env
	@$(terraform) plan $(terraform_plan_args) $(terraform_vars) -out plan

plan_destroy: $(terraform) check-deploy-env
	@$(terraform) plan --destroy $(terraform_plan_args) $(terraform_vars) -out plan

plan-artifacts.tar:
	tar $(tar_ignore_arg) $(tar_absolute_names) -cvf $@ \
		plan .terraform.lock.hcl

refresh: $(terraform) check-deploy-env
	$(terraform) refresh $(terraform_vars)

apply: $(terraform) check-deploy-env
	$(terraform) apply -auto-approve plan

lint_check: $(terraform)
	$(terraform) fmt -check -recursive .

lint_fix: $(terraform)
	$(terraform) fmt -recursive .
