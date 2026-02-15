variable "docker_image" {
  description = "The full ECR image URI to deploy"
  type        = string
  default     = null
}

variable "oidc_client_secret" {
  description = "Google OAuth 2.0 Client Secret for ALB OIDC authentication"
  type        = string
  sensitive   = true
}
