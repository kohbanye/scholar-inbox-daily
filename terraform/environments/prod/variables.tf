variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "ap-northeast-1"
}

variable "ecr_repository_name" {
  description = "Name of the ECR repository"
  type        = string
  default     = "scholar-inbox-daily"
}

variable "lambda_function_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "scholar-inbox-daily"
}

variable "scholar_inbox_email" {
  description = "Email for Scholar Inbox authentication"
  type        = string
  sensitive   = true
}

variable "scholar_inbox_password" {
  description = "Password for Scholar Inbox authentication"
  type        = string
  sensitive   = true
}

variable "slack_api_token" {
  description = "Slack API token for posting messages"
  type        = string
  sensitive   = true
}

variable "slack_channel_id" {
  description = "Slack channel ID to post messages to"
  type        = string
  sensitive   = true
} 