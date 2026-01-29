output "bucket_name" {
  description = "S3 bucket for archived data"
  value       = aws_s3_bucket.archive.id
}

output "bucket_arn" {
  description = "S3 bucket ARN"
  value       = aws_s3_bucket.archive.arn
}

output "lambda_function_name" {
  description = "Lambda function name"
  value       = aws_lambda_function.archiver.function_name
}

output "lambda_function_arn" {
  description = "Lambda function ARN"
  value       = aws_lambda_function.archiver.arn
}

output "eventbridge_rule_arn" {
  description = "EventBridge schedule rule ARN"
  value       = aws_cloudwatch_event_rule.schedule.arn
}

output "log_group_name" {
  description = "CloudWatch Log Group for Lambda"
  value       = aws_cloudwatch_log_group.lambda.name
}
