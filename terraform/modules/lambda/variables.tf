variable "ecr_repository_name" {
  description = "Name of the ECR repository"
  type        = string
}

variable "lambda_function_name" {
  description = "Name of the Lambda function"
  type        = string
}

variable "environment_variables" {
  description = "Environment variables for the Lambda function"
  type        = map(string)
  sensitive   = true
}

variable "image_tag" {
  description = "The tag of the Docker image to use for the Lambda function"
  type        = string
}
