resource "google_project_service" "artifactregistry" {
  project = local.repo_project
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "sql_admin" {
  project = local.repo_project
  service = "sqladmin.googleapis.com"
}

resource "google_artifact_registry_repository" "default" {
  depends_on             = [google_project_service.artifactregistry]
  location               = local.repo_region
  repository_id          = "knz-myservice"
  description            = "Docker Repository for My Service Service"
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "delete-old-untagged"
    action = "DELETE"
    condition {
      tag_state  = "UNTAGGED"
      older_than = "86400s"
    }
  }
}

data "google_project" "env" {
  for_each   = toset(local.env_projects)
  project_id = each.key
}

resource "google_project_iam_member" "cloudrun_sa_registry_reader" {
  for_each = data.google_project.env

  project = local.repo_project
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:service-${each.value.number}@serverless-robot-prod.iam.gserviceaccount.com"
}
