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
env_root=environments/$(DEPLOY_ENV)
terraform_backend_prefix_file=state-prefix.s3.tfbackend
terraform_args=-chdir="$(env_root)"
terraform_vars=-var "app_name=$(app_name)"
terraform_plan_args=-input=false

ifdef CI
terraform=$(shell which terraform)
terraform_plan_args+=-no-color
else
terraform=$(bin)/terraform
endif
terraform_with_args=$(terraform) $(terraform_args)

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

%/$(terraform_backend_prefix_file):
	echo 'key  = "$(app_name)-state"' > $@

all_tf_lock_files := $(patsubst %/versions.tf,%/.terraform.lock.hcl,$(wildcard environments/*/versions.tf))

environments/%/.terraform.lock.hcl: environments/%/versions.tf environments/%/providers.tf
	$(terraform) -chdir=$(@D) init -backend=false -upgrade

.PHONY: providers_lock
providers_lock: $(terraform) $(all_tf_lock_files)

init: $(terraform) check-deploy-env $(env_root)/$(terraform_backend_prefix_file)
	$(terraform_with_args) init -backend-config=$(terraform_backend_prefix_file)

init_upgrade: $(terraform) check-deploy-env $(env_root)/$(terraform_backend_prefix_file)
	$(terraform_with_args) init -backend-config=$(terraform_backend_prefix_file) -upgrade

import: $(terraform) check-deploy-env
	@echo $(terraform_with_args) import $(terraform_vars)

state: $(terraform) check-deploy-env
	@echo $(terraform_with_args) state

plan: $(terraform) check-deploy-env
	@$(terraform_with_args) plan $(terraform_plan_args) $(terraform_vars) -out plan

plan_destroy: $(terraform) check-deploy-env
	@$(terraform_with_args) plan --destroy $(terraform_plan_args) $(terraform_vars) -out plan

plan-artifacts.tar:
	tar $(tar_ignore_arg) $(tar_absolute_names) -cvf $@ \
		plan .terraform.lock.hcl

refresh: $(terraform) check-deploy-env
	$(terraform_with_args) refresh $(terraform_vars)

apply: $(terraform) check-deploy-env
	$(terraform_with_args) apply -auto-approve plan

lint_check: $(terraform)
	$(terraform) fmt -check -recursive .

lint_fix: $(terraform)
	$(terraform) fmt -recursive .
