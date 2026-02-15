resource "google_compute_region_network_endpoint_group" "myservice" {
  name                  = "myservice-neg"
  project               = var.project
  network_endpoint_type = "SERVERLESS"
  region                = var.region
  cloud_run {
    service = "myservice"
  }
}

resource "google_compute_region_backend_service" "external" {
  name                  = "external"
  description           = "This backend service is accessible by anyone online"
  project               = var.project
  region                = var.region
  load_balancing_scheme = "EXTERNAL_MANAGED"
  protocol              = "HTTPS"
  log_config {
    enable = true
  }
  backend {
    group           = google_compute_region_network_endpoint_group.myservice.self_link
    capacity_scaler = 1
    balancing_mode  = "UTILIZATION"
  }
}

resource "google_compute_region_backend_service" "internal" {
  name                  = "internal"
  description           = "This backend service is only accessible by internal personnel who are part of Kynzy's google organization"
  project               = var.project
  region                = var.region
  load_balancing_scheme = "EXTERNAL_MANAGED"
  protocol              = "HTTPS"
  iap {
    oauth2_client_id     = google_iap_client.project_client.client_id
    oauth2_client_secret = google_iap_client.project_client.secret
  }
  log_config {
    enable = true
  }
  backend {
    group           = google_compute_region_network_endpoint_group.myservice.self_link
    capacity_scaler = 1
    balancing_mode  = "UTILIZATION"
  }
}
