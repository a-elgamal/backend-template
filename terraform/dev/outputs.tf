output "uri" {
  description = "The URI of the deployed myservice service"
  value       = module.myservice.uri
}

output "internal-backend-service" {
  description = "The URI for the myservice service intenral backend that can be used with a regional load balancer"
  value       = module.myservice.internal-backend-service
}

output "external-backend-service" {
  description = "The URI for the myservice service external backend that can be used with a regional load balancer"
  value       = module.myservice.external-backend-service
}
