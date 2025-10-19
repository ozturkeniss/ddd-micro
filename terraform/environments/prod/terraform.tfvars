# Production Environment Configuration

environment = "prod"
aws_region  = "us-west-2"

# VPC Configuration
vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]

# EKS Configuration (larger instances for prod)
eks_node_groups = {
  general = {
    instance_types = ["t3.large", "t3.xlarge"]
    capacity_type  = "ON_DEMAND"
    min_size      = 2
    max_size      = 5
    desired_size  = 3
    disk_size     = 50
  }
  spot = {
    instance_types = ["t3.large", "t3.xlarge", "m5.large", "m5.xlarge"]
    capacity_type  = "SPOT"
    min_size      = 0
    max_size      = 10
    desired_size  = 2
    disk_size     = 50
  }
}

# RDS Configuration (larger instances for prod)
rds_databases = {
  user = {
    engine_version = "15.4"
    instance_class = "db.r5.large"
    allocated_storage = 100
    max_allocated_storage = 1000
    storage_encrypted = true
    backup_retention_period = 30
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = true
    deletion_protection = true
  }
  product = {
    engine_version = "15.4"
    instance_class = "db.r5.large"
    allocated_storage = 100
    max_allocated_storage = 1000
    storage_encrypted = true
    backup_retention_period = 30
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = true
    deletion_protection = true
  }
  basket = {
    engine_version = "15.4"
    instance_class = "db.r5.large"
    allocated_storage = 100
    max_allocated_storage = 1000
    storage_encrypted = true
    backup_retention_period = 30
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = true
    deletion_protection = true
  }
  payment = {
    engine_version = "15.4"
    instance_class = "db.r5.large"
    allocated_storage = 100
    max_allocated_storage = 1000
    storage_encrypted = true
    backup_retention_period = 30
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = true
    deletion_protection = true
  }
}

# Redis Configuration (larger for prod)
redis_config = {
  node_type = "cache.r5.large"
  num_cache_nodes = 2
  parameter_group_name = "default.redis7"
  engine_version = "7.0"
  port = 6379
  maintenance_window = "sun:05:00-sun:06:00"
  snapshot_retention_limit = 30
  snapshot_window = "03:00-05:00"
}

# Kafka Configuration (larger for prod)
kafka_config = {
  kafka_version = "3.4.0"
  number_of_broker_nodes = 3
  broker_node_group_info = {
    instance_type = "kafka.r5.large"
    ebs_volume_size = 100
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
        log_group = "/aws/msk/ddd-micro-prod"
      }
    }
  }
}
