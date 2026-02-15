variable "env" {
  description = "The name of the environment that the service will be deployed to"
  type        = string
  nullable    = false
}

variable "project" {
  description = "The service project id to deploy the service to"
  type        = string
  nullable    = false
}

variable "host_project" {
  description = "The host project id for the environment"
  type        = string
  nullable    = false
}

variable "region" {
  description = "The region to deploy the service to"
  type        = string
  nullable    = false
}

variable "repo_region" {
  description = "The region in which artifact registry has the image service region"
  type        = string
  default     = "europe-west2"
  nullable    = false
}

variable "port" {
  description = "The port that the service should be listening to"
  type        = number
  default     = 8080
  nullable    = false
}

variable "docker_image_tag" {
  description = "The tag of the docker image to deploy"
  type        = string
  default     = "latest"
  nullable    = false
}

variable "docker_image_digest" {
  description = "The digest of the docker image to deploy"
  type        = string
  nullable    = true
}

variable "cloudsql_instance_type" {
  description = "The type of CloudSQL instance to provision"
  type        = string
  default     = "db-f1-micro"
}

variable "cloudsql_ha_enabled" {
  description = "Whether High Availability should be enabled or not on the CloudSQL instance"
  type        = bool
  default     = true
}
