# Development Environment Configuration

environment = "dev"
aws_region  = "us-west-2"

# VPC Configuration
vpc_cidr = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]

# EKS Configuration
eks_node_groups = {
  general = {
    instance_types = ["t3.medium"]
    capacity_type  = "ON_DEMAND"
    min_size      = 1
    max_size      = 2
    desired_size  = 1
    disk_size     = 20
  }
}

# RDS Configuration (smaller instances for dev)
rds_databases = {
  user = {
    engine_version = "15.4"
    instance_class = "db.t3.micro"
    allocated_storage = 20
    max_allocated_storage = 50
    storage_encrypted = false
    backup_retention_period = 1
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = false
    deletion_protection = false
  }
  product = {
    engine_version = "15.4"
    instance_class = "db.t3.micro"
    allocated_storage = 20
    max_allocated_storage = 50
    storage_encrypted = false
    backup_retention_period = 1
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = false
    deletion_protection = false
  }
  basket = {
    engine_version = "15.4"
    instance_class = "db.t3.micro"
    allocated_storage = 20
    max_allocated_storage = 50
    storage_encrypted = false
    backup_retention_period = 1
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = false
    deletion_protection = false
  }
  payment = {
    engine_version = "15.4"
    instance_class = "db.t3.micro"
    allocated_storage = 20
    max_allocated_storage = 50
    storage_encrypted = false
    backup_retention_period = 1
    backup_window = "03:00-04:00"
    maintenance_window = "sun:04:00-sun:05:00"
    multi_az = false
    deletion_protection = false
  }
}

# Redis Configuration (smaller for dev)
redis_config = {
  node_type = "cache.t3.micro"
  num_cache_nodes = 1
  parameter_group_name = "default.redis7"
  engine_version = "7.0"
  port = 6379
  maintenance_window = "sun:05:00-sun:06:00"
  snapshot_retention_limit = 1
  snapshot_window = "03:00-05:00"
}

# Kafka Configuration (smaller for dev)
kafka_config = {
  kafka_version = "3.4.0"
  number_of_broker_nodes = 1
  broker_node_group_info = {
    instance_type = "kafka.t3.small"
    ebs_volume_size = 10
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
        log_group = "/aws/msk/ddd-micro-dev"
      }
    }
  }
}
