# DDD Microservices Infrastructure Outputs

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = module.vpc.private_subnet_ids
}

# EKS Outputs
output "eks_cluster_id" {
  description = "ID of the EKS cluster"
  value       = module.eks.cluster_id
}

output "eks_cluster_arn" {
  description = "ARN of the EKS cluster"
  value       = module.eks.cluster_arn
}

output "eks_cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "eks_cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster"
  value       = module.eks.cluster_security_group_id
}

output "eks_oidc_provider_arn" {
  description = "ARN of the OIDC Provider"
  value       = module.eks.oidc_provider_arn
}

# RDS Outputs
output "rds_endpoints" {
  description = "RDS instance endpoints"
  value       = module.rds.endpoints
}

output "rds_security_group_id" {
  description = "Security group ID for RDS instances"
  value       = module.rds.security_group_id
}

# Redis Outputs
output "redis_endpoint" {
  description = "Redis cluster endpoint"
  value       = module.redis.endpoint
}

output "redis_port" {
  description = "Redis port"
  value       = module.redis.port
}

# ECR Outputs
output "ecr_repository_urls" {
  description = "ECR repository URLs"
  value       = module.ecr.repository_urls
}

# MSK Outputs
output "msk_cluster_arn" {
  description = "ARN of the MSK cluster"
  value       = module.msk.cluster_arn
}

output "msk_bootstrap_brokers" {
  description = "Plaintext connection host:port pairs"
  value       = module.msk.bootstrap_brokers
}

output "msk_bootstrap_brokers_tls" {
  description = "TLS connection host:port pairs"
  value       = module.msk.bootstrap_brokers_tls
}

# Common Outputs
output "environment" {
  description = "Environment name"
  value       = var.environment
}

output "aws_region" {
  description = "AWS region"
  value       = var.aws_region
}
