resource "google_project_service" "monitoring" {
  project = var.project
  service = "monitoring.googleapis.com"
}

resource "google_project_service" "cloud_trace" {
  project = var.project
  service = "cloudtrace.googleapis.com"
}

resource "google_project_service" "logging" {
  project = var.project
  service = "logging.googleapis.com"
}

resource "google_project_service" "secret_manager" {
  project = var.project
  service = "secretmanager.googleapis.com"
}

resource "google_project_service" "sql_admin" {
  project = var.project
  service = "sqladmin.googleapis.com"
}

resource "google_project_service" "iap" {
  project = var.project
  service = "iap.googleapis.com"
}
