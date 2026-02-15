output "uri" {
  description = "The URI of the deployed myservice service"
  value       = google_cloud_run_v2_service.default.uri
}

output "internal-backend-service" {
  description = "The URI for the myservice service intenral backend that can be used with a regional load balancer"
  value       = google_compute_region_backend_service.internal.self_link
}

output "external-backend-service" {
  description = "The URI for the myservice service external backend that can be used with a regional load balancer"
  value       = google_compute_region_backend_service.external.self_link
}
