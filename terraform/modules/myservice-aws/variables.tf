variable "env" {
  description = "The name of the environment (e.g. dev, staging, prod)"
  type        = string
  nullable    = false
}

variable "region" {
  description = "The AWS region to deploy to"
  type        = string
  nullable    = false
}

variable "vpc_id" {
  description = "The VPC ID to deploy into"
  type        = string
  nullable    = false
}

variable "public_subnet_ids" {
  description = "Subnet IDs for the ALB (must be public subnets in at least 2 AZs)"
  type        = list(string)
  nullable    = false
}

variable "private_subnet_ids" {
  description = "Subnet IDs for ECS tasks and RDS (private subnets)"
  type        = list(string)
  nullable    = false
}

variable "port" {
  description = "The port that the service listens on"
  type        = number
  default     = 8080
  nullable    = false
}

variable "docker_image" {
  description = "The full ECR image URI including tag or digest (e.g. 123456789.dkr.ecr.eu-west-1.amazonaws.com/myservice:latest)"
  type        = string
  nullable    = false
}

variable "oidc_client_id" {
  description = "Google OAuth 2.0 Client ID for ALB OIDC authentication"
  type        = string
  nullable    = false
}

variable "oidc_client_secret" {
  description = "Google OAuth 2.0 Client Secret for ALB OIDC authentication"
  type        = string
  sensitive   = true
  nullable    = false
}

variable "rds_instance_class" {
  description = "The RDS instance class"
  type        = string
  default     = "db.t4g.micro"
}

variable "rds_multi_az" {
  description = "Whether to enable Multi-AZ for RDS"
  type        = bool
  default     = false
}

variable "certificate_arn" {
  description = "ACM certificate ARN for the ALB HTTPS listener"
  type        = string
  nullable    = false
}

variable "task_cpu" {
  description = "CPU units for the ECS task (1024 = 1 vCPU)"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "Memory (MiB) for the ECS task"
  type        = number
  default     = 512
}
