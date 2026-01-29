variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-2"
}

variable "aws_profile" {
  description = "AWS CLI profile to use"
  type        = string
  default     = "personal"
}

variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
  default     = "517397653073"
}

variable "bucket_name" {
  description = "S3 bucket name for archived data"
  type        = string
  default     = "nes-outage-status-checker-archive"
}

variable "retention_days" {
  description = "Number of days to retain archived data"
  type        = number
  default     = 90
}
