resource "google_service_account" "default" {
  account_id   = "app-server"
  display_name = "The Service Account used by myservice app server"
  project      = var.project
}

resource "google_project_iam_member" "sa_metric_writer" {
  project = var.project
  role    = "roles/monitoring.metricWriter"
  member  = google_service_account.default.member
}

resource "google_project_iam_member" "sa_log_writer" {
  project = var.project
  role    = "roles/logging.logWriter"
  member  = google_service_account.default.member
}

resource "google_project_iam_member" "sa_cloud_trace_agent" {
  project = var.project
  role    = "roles/cloudtrace.agent"
  member  = google_service_account.default.member
}

resource "google_project_iam_member" "sa_cloudsql_client" {
  project = var.project
  role    = "roles/cloudsql.client"
  member  = google_service_account.default.member
}

resource "google_project_iam_member" "sa_cloudsql_user" {
  project = var.project
  role    = "roles/cloudsql.user"
  member  = google_service_account.default.member
}

resource "google_secret_manager_secret_iam_member" "sa_secret_access" {
  secret_id = google_secret_manager_secret.otlp_collector_config.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.default.member
  project   = var.project
}

resource "google_service_account_iam_member" "repo_sa" {
  service_account_id = google_service_account.default.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:github-actions@knz-myservice-repo.iam.gserviceaccount.com"
}
