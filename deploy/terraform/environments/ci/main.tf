module "test_api_infra" {
  source                = "../../modules/app-infra"
  resources_prefix      = "${var.resources_prefix}test-"
  resources_description = var.resources_description
  local_stack_env       = true
}
module "test_api_service" {
  source                   = "../../modules/api-service"
  app_name                 = var.app_name
  app_bus_name             = module.test_api_infra.app_bus_name
  dummy_messages_topic_arn = module.test_api_infra.dummy_messages_topic_arn
  resources_prefix         = "${var.resources_prefix}test-"
  resources_description    = var.resources_description
  aws_primary_region       = var.aws_primary_region
  service_endpoint         = "http://host.docker.internal:48080"
}