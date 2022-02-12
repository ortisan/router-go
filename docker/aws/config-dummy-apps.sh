#!/bin/bash
set -x
# Create parameter App 1
awslocal ssm put-parameter --name "app1.error.threshold" --type String --value "0" --overwrite --region=sa-east-1
# Create parameter App 2
awslocal ssm put-parameter --name "app2.error.threshold" --type String --value "50" --overwrite --region=sa-east-1
# Create parameter App 3
awslocal ssm put-parameter --name "app3.error.threshold" --type String --value "100" --overwrite --region=sa-east-1
set +x
