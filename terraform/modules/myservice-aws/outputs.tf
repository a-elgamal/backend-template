output "alb_dns_name" {
  description = "The DNS name of the ALB"
  value       = aws_lb.default.dns_name
}

output "ecs_service_arn" {
  description = "The ARN of the ECS service"
  value       = aws_ecs_service.default.id
}

output "ecs_cluster_name" {
  description = "The name of the ECS cluster"
  value       = aws_ecs_cluster.default.name
}

output "rds_endpoint" {
  description = "The RDS instance endpoint"
  value       = aws_db_instance.default.endpoint
}
