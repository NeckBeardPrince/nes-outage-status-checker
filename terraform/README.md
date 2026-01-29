# NES Outage Archiver - Terraform

Infrastructure as Code for deploying the NES Outage Data Archiver to AWS.

## Resources Created

- **S3 Bucket** - `nes-outage-status-checker-archive` for storing archived data
- **Lambda Function** - `nes-outage-archiver` to fetch and store data
- **IAM Role & Policies** - For Lambda to write to S3 and CloudWatch Logs
- **EventBridge Rule** - Triggers Lambda every 10 minutes
- **CloudWatch Log Group** - 14-day retention for Lambda logs

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- AWS CLI configured with `personal` profile
- AWS credentials with permissions to create the above resources

## Quick Start

```bash
cd terraform

# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Deploy
terraform apply
```

## Configuration

Default values are set in `variables.tf`. Override as needed:

```bash
# Use different values
terraform apply \
  -var="bucket_name=my-custom-bucket" \
  -var="retention_days=30"
```

| Variable | Default | Description |
|----------|---------|-------------|
| `aws_region` | `us-east-2` | AWS region |
| `aws_profile` | `personal` | AWS CLI profile |
| `aws_account_id` | `517397653073` | AWS Account ID |
| `bucket_name` | `nes-outage-status-checker-archive` | S3 bucket name |
| `retention_days` | `90` | Days to keep archived data |

## Outputs

After applying, Terraform will output:

- `bucket_name` - S3 bucket for archived data
- `lambda_function_name` - Lambda function name
- `lambda_function_arn` - Lambda function ARN
- `eventbridge_rule_arn` - Schedule rule ARN
- `log_group_name` - CloudWatch Log Group

## Testing

Manually invoke the Lambda to test:

```bash
aws lambda invoke \
  --profile personal \
  --region us-east-2 \
  --function-name nes-outage-archiver \
  --payload '{}' \
  response.json

cat response.json
```

Check the S3 bucket:

```bash
aws s3 ls s3://nes-outage-status-checker-archive/ --profile personal --recursive
```

## Tear Down

```bash
# Empty the bucket first (required before deletion)
aws s3 rm s3://nes-outage-status-checker-archive --recursive --profile personal

# Destroy all resources
terraform destroy
```

## Cost Estimate

| Resource | Monthly Cost |
|----------|-------------|
| Lambda (4,320 invocations) | ~$0.01 |
| S3 (~180MB storage) | ~$0.01 |
| EventBridge | Free |
| CloudWatch Logs | ~$0.01 |
| **Total** | **< $0.05/month** |
