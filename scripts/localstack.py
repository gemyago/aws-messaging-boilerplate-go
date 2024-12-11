#!/usr/bin/env python

import argparse
import boto3

def create_queue(queue_name, region='us-east-2'):
    sqs = boto3.client(
        'sqs',
        endpoint_url='http://localhost:4566',
        region_name=region,
        aws_access_key_id='dummy',
        aws_secret_access_key='dummy'
    )
    response = sqs.create_queue(QueueName=queue_name)
    print(f"Queue {queue_name} created with URL: {response['QueueUrl']}")

def provision_queues(region):
    create_queue('messages-queue', region)
    create_queue('test-messages-queue', region)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='LocalStack SQS Provisioning Script')
    parser.add_argument('command', choices=['provision'], help='Command to execute')
    args = parser.parse_args()

    if not args.command:
        parser.print_help()
    else:
        if args.command == 'provision':
            provision_queues(region='us-east-1')
