output "queues" {
  value = [
    for index, subscription in local.sqs_subscriptions : {
      subscription_id = subscription.key
      url             = aws_sqs_queue.primary[subscription.key].url
      dlq_url         = aws_sqs_queue.dead_letter[subscription.key].url
    }
  ]
}