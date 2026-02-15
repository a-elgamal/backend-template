resource "google_secret_manager_secret" "otlp_collector_config" {
  depends_on = [google_project_service.secret_manager]
  secret_id  = "otlp_collector_config"
  project    = var.project
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "otlp_collector_config" {
  secret      = google_secret_manager_secret.otlp_collector_config.name
  secret_data = <<-EOT
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
  syslog:
    tcp:
      listen_address: 0.0.0.0:1601
    protocol: rfc3164

processors:
  batch:
  memory_limiter:
    check_interval: 1s
    limit_percentage: 65
    spike_limit_percentage: 20
  resourcedetection:
    detectors: [env, gcp]
    timeout: 10s

exporters:
  googlecloud:
    log:
      default_log_name: opentelemetry.io/collector-exported-log

extensions:
  health_check:
    endpoint: 0.0.0.0:13133

service:
  extensions: [health_check]
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [googlecloud]
    traces:
      receivers: [otlp]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [googlecloud]
    logs:
      receivers: [otlp, syslog]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [googlecloud]
  EOT
}
