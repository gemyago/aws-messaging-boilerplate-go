module "app" {
  source                = "./app"
  app_name              = var.app_name
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description
  aws_primary_region    = var.aws_primary_region
  local_stack_env       = var.local_stack_env
}

module "test_app" {
  count                 = var.setup_test_env ? 1 : 0
  source                = "./app"
  app_name              = var.app_name
  resources_prefix      = "${var.resources_prefix}test-"
  resources_description = var.resources_description
  aws_primary_region    = var.aws_primary_region
  local_stack_env       = var.local_stack_env
}