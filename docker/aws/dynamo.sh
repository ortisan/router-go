#!/bin/bash
set -x
# Cria o bucket
awslocal dynamodb --endpoint-url=http://localstack:4566 create-table \
    --table-name HealthCheck \
    --attribute-definitions \
        AttributeName=ID,AttributeType=S \
        AttributeName=ServiceURL,AttributeType=S \
        AttributeName=Value,AttributeType=N \
        AttributeName=EpochTimestamp,AttributeType=N \
    --key-schema \
        AttributeName=ID,KeyType=HASH \
        AttributeName=EpochTimestamp,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=10

set +x