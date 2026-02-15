resource "google_project_service_identity" "iap_sa" {
  provider = google-beta

  project = var.project
  service = "iap.googleapis.com"
}

resource "google_project_iam_member" "iap_sa_run_invoker" {
  project = var.project
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_project_service_identity.iap_sa.email}"
}

resource "google_iap_brand" "project_brand" {
  depends_on = [google_project_service.iap]

  support_email     = "TODO_SUPPORT_EMAIL"
  application_title = "My Service"
  project           = var.project
}

resource "google_iap_client" "project_client" {
  display_name = "My Service Internal APIs"
  brand        = google_iap_brand.project_brand.name
}

resource "google_iap_web_region_backend_service_iam_binding" "admins" {
  project                    = var.project
  region                     = var.region
  web_region_backend_service = google_compute_region_backend_service.internal.name
  role                       = "roles/iap.httpsResourceAccessor"
  members = [
    "user:TODO_ADMIN_EMAIL",
  ]
}
