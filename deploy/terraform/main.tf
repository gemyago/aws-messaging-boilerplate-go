module "app" {
  source                = "./app"
  app_name              = var.app_name
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description
  aws_primary_region    = var.aws_primary_region
  local_stack_env       = var.local_stack_env
  service_endpoint      = var.service_endpoint
}