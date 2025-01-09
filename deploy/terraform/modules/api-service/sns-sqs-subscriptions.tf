module "sns_sqs_subscriptions" {
  source                = "../sns-sqs-subscriptions"
  app_name              = var.app_name
  resources_prefix      = var.resources_prefix
  resources_description = var.resources_description
  subscriptions = [
    {
      subscription_id = "dummy-messages"
      topic_arn       = var.dummy_messages_topic_arn
    }
  ]
}