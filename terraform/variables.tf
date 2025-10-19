# DDD Microservices Infrastructure Variables

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "ddd-micro"
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]
}

# EKS Configuration
variable "eks_node_groups" {
  description = "EKS node groups configuration"
  type = map(object({
    instance_types = list(string)
    capacity_type  = string
    min_size      = number
    max_size      = number
    desired_size  = number
    disk_size     = number
  }))
  default = {
    general = {
      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"
      min_size      = 1
      max_size      = 3
      desired_size  = 2
      disk_size     = 20
    }
    spot = {
      instance_types = ["t3.medium", "t3.large"]
      capacity_type  = "SPOT"
      min_size      = 0
      max_size      = 5
      desired_size  = 1
      disk_size     = 20
    }
  }
}

# RDS Configuration
variable "rds_databases" {
  description = "RDS databases configuration"
  type = map(object({
    engine_version = string
    instance_class = string
    allocated_storage = number
    max_allocated_storage = number
    storage_encrypted = bool
    backup_retention_period = number
    backup_window = string
    maintenance_window = string
    multi_az = bool
    deletion_protection = bool
  }))
  default = {
    user = {
      engine_version = "15.4"
      instance_class = "db.t3.micro"
      allocated_storage = 20
      max_allocated_storage = 100
      storage_encrypted = true
      backup_retention_period = 7
      backup_window = "03:00-04:00"
      maintenance_window = "sun:04:00-sun:05:00"
      multi_az = false
      deletion_protection = false
    }
    product = {
      engine_version = "15.4"
      instance_class = "db.t3.micro"
      allocated_storage = 20
      max_allocated_storage = 100
      storage_encrypted = true
      backup_retention_period = 7
      backup_window = "03:00-04:00"
      maintenance_window = "sun:04:00-sun:05:00"
      multi_az = false
      deletion_protection = false
    }
    basket = {
      engine_version = "15.4"
      instance_class = "db.t3.micro"
      allocated_storage = 20
      max_allocated_storage = 100
      storage_encrypted = true
      backup_retention_period = 7
      backup_window = "03:00-04:00"
      maintenance_window = "sun:04:00-sun:05:00"
      multi_az = false
      deletion_protection = false
    }
    payment = {
      engine_version = "15.4"
      instance_class = "db.t3.micro"
      allocated_storage = 20
      max_allocated_storage = 100
      storage_encrypted = true
      backup_retention_period = 7
      backup_window = "03:00-04:00"
      maintenance_window = "sun:04:00-sun:05:00"
      multi_az = false
      deletion_protection = false
    }
  }
}

# Redis Configuration
variable "redis_config" {
  description = "Redis configuration"
  type = object({
    node_type = string
    num_cache_nodes = number
    parameter_group_name = string
    engine_version = string
    port = number
    maintenance_window = string
    snapshot_retention_limit = number
    snapshot_window = string
  })
  default = {
    node_type = "cache.t3.micro"
    num_cache_nodes = 1
    parameter_group_name = "default.redis7"
    engine_version = "7.0"
    port = 6379
    maintenance_window = "sun:05:00-sun:06:00"
    snapshot_retention_limit = 5
    snapshot_window = "03:00-05:00"
  }
}

# ECR Configuration
variable "ecr_repositories" {
  description = "ECR repositories"
  type        = list(string)
  default     = [
    "user-service",
    "product-service", 
    "basket-service",
    "payment-service",
    "krakend-gateway"
  ]
}

# Kafka Configuration
variable "kafka_config" {
  description = "MSK Kafka configuration"
  type = object({
    kafka_version = string
    number_of_broker_nodes = number
    broker_node_group_info = object({
      instance_type = string
      ebs_volume_size = number
    })
    encryption_info = object({
      encryption_at_rest_kms_key_id = string
      encryption_in_transit = object({
        client_broker = string
        in_cluster = bool
      })
    })
    logging_info = object({
      broker_logs = object({
        cloudwatch_logs = object({
          enabled = bool
          log_group = string
        })
      })
    })
  })
  default = {
    kafka_version = "3.4.0"
    number_of_broker_nodes = 2
    broker_node_group_info = {
      instance_type = "kafka.t3.small"
      ebs_volume_size = 20
    }
    encryption_info = {
      encryption_at_rest_kms_key_id = ""
      encryption_in_transit = {
        client_broker = "TLS"
        in_cluster = true
      }
    }
    logging_info = {
      broker_logs = {
        cloudwatch_logs = {
          enabled = true
          log_group = "/aws/msk/ddd-micro"
        }
      }
    }
  }
}
