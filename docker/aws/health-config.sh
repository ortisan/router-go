#!/bin/bash
set -x
# Creat sns topic
awslocal sns create-topic --name health-cells-topic --endpoint-url=http://localstack:4566 --region=sa-east-1
# Create sqs queue
awslocal sqs create-queue --queue-name health-cells-queue --endpoint-url=http://localstack:4566 --region=sa-east-1
# Subscribe topic with sqs
awslocal sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:health-cells-topic --protocol sqs --notification-endpoint http://localstack:4566/000000000000/health-cells-queue
set +x
