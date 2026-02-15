data "google_compute_network" "default" {
  project = var.host_project
  name    = var.env
}

data "google_compute_subnetwork" "default" {
  project = var.host_project
  name    = "myservice-${trimspace(lower(var.env))}-${trimspace(lower(var.region))}"
  region  = var.region
}
