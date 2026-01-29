terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
  }
}

provider "aws" {
  region  = var.aws_region
  profile = var.aws_profile
}

# -----------------------------------------------------------------------------
# S3 Bucket for archived data
# -----------------------------------------------------------------------------
resource "aws_s3_bucket" "archive" {
  bucket = var.bucket_name

  tags = {
    Name        = "NES Outage Archive"
    Project     = "nes-outage-status-checker"
    ManagedBy   = "terraform"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "archive" {
  bucket = aws_s3_bucket.archive.id

  rule {
    id     = "expire-old-data"
    status = "Enabled"

    expiration {
      days = var.retention_days
    }
  }
}

resource "aws_s3_bucket_public_access_block" "archive" {
  bucket = aws_s3_bucket.archive.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "archive" {
  bucket = aws_s3_bucket.archive.id

  versioning_configuration {
    status = "Disabled"
  }
}

# -----------------------------------------------------------------------------
# IAM Role for Lambda
# -----------------------------------------------------------------------------
resource "aws_iam_role" "lambda" {
  name = "nes-outage-archiver-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Project   = "nes-outage-status-checker"
    ManagedBy = "terraform"
  }
}

# Policy for CloudWatch Logs
resource "aws_iam_role_policy" "lambda_logs" {
  name = "cloudwatch-logs"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:*"
      }
    ]
  })
}

# Policy for S3 access
resource "aws_iam_role_policy" "lambda_s3" {
  name = "s3-write"
  role = aws_iam_role.lambda.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:PutObjectAcl"
        ]
        Resource = "${aws_s3_bucket.archive.arn}/*"
      }
    ]
  })
}

# -----------------------------------------------------------------------------
# Lambda Function
# -----------------------------------------------------------------------------
data "archive_file" "lambda" {
  type        = "zip"
  source_file = "${path.module}/../lambda/nes_archiver.py"
  output_path = "${path.module}/lambda.zip"
}

resource "aws_lambda_function" "archiver" {
  function_name = "nes-outage-archiver"
  description   = "Fetches NES outage data and archives to S3 every 10 minutes"

  filename         = data.archive_file.lambda.output_path
  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "python3.11"
  handler = "nes_archiver.lambda_handler"
  timeout = 30
  memory_size = 128

  role = aws_iam_role.lambda.arn

  environment {
    variables = {
      BUCKET_NAME = aws_s3_bucket.archive.id
    }
  }

  tags = {
    Project   = "nes-outage-status-checker"
    ManagedBy = "terraform"
  }
}

# CloudWatch Log Group for Lambda
resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${aws_lambda_function.archiver.function_name}"
  retention_in_days = 14

  tags = {
    Project   = "nes-outage-status-checker"
    ManagedBy = "terraform"
  }
}

# -----------------------------------------------------------------------------
# EventBridge Rule (CloudWatch Events)
# -----------------------------------------------------------------------------
resource "aws_cloudwatch_event_rule" "schedule" {
  name                = "nes-outage-archiver-schedule"
  description         = "Trigger NES outage archiver every 10 minutes"
  schedule_expression = "rate(10 minutes)"

  tags = {
    Project   = "nes-outage-status-checker"
    ManagedBy = "terraform"
  }
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule      = aws_cloudwatch_event_rule.schedule.name
  target_id = "nes-outage-archiver"
  arn       = aws_lambda_function.archiver.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.archiver.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.schedule.arn
}
