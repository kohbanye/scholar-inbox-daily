provider "aws" {
  region = var.aws_region
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "scholar-inbox-daily-terraform-state"
    key    = "prod/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

module "lambda" {
  source = "../../modules/lambda"

  ecr_repository_name  = var.ecr_repository_name
  lambda_function_name = var.lambda_function_name
  environment_variables = {
    SCHOLAR_INBOX_EMAIL    = var.scholar_inbox_email
    SCHOLAR_INBOX_PASSWORD = var.scholar_inbox_password
    SLACK_API_TOKEN        = var.slack_api_token
    SLACK_CHANNEL_ID       = var.slack_channel_id
    TZ                     = "Asia/Tokyo"
  }
} 