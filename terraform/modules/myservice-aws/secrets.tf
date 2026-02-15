resource "aws_secretsmanager_secret" "otlp_collector_config" {
  name = "myservice-${var.env}-otlp-collector-config"
}

resource "aws_secretsmanager_secret_version" "otlp_collector_config" {
  secret_id = aws_secretsmanager_secret.otlp_collector_config.id
  secret_string = <<-EOT
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
    detectors: [env, ecs]
    timeout: 10s

exporters:
  awsxray:
  awsemf:
    namespace: myservice/${var.env}
    log_group_name: /ecs/myservice-${var.env}

extensions:
  health_check:
    endpoint: 0.0.0.0:13133

service:
  extensions: [health_check]
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [awsemf]
    traces:
      receivers: [otlp]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [awsxray]
    logs:
      receivers: [otlp, syslog]
      processors: [batch, memory_limiter, resourcedetection]
      exporters: [awsemf]
  EOT
}
