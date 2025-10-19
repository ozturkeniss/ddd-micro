# DDD Microservices Architecture

A comprehensive microservices architecture built with Domain-Driven Design principles, featuring user management, product catalog, shopping basket, and payment processing services.

## Architecture Overview

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Client Layer"
        WEB[Web Client]
        MOBILE[Mobile Client]
    end
    
    subgraph "API Gateway"
        KRAKEND[KrakenD Gateway]
    end
    
    subgraph "Microservices"
        USER[User Service<br/>Port: 8080<br/>gRPC: 9091]
        PRODUCT[Product Service<br/>Port: 8081<br/>gRPC: 9092]
        BASKET[Basket Service<br/>Port: 8083<br/>gRPC: 9093]
        PAYMENT[Payment Service<br/>Port: 8084]
    end
    
    subgraph "Message Broker"
        KAFKA[Apache Kafka<br/>Port: 9092]
    end
    
    subgraph "Data Layer"
        USER_DB[(User DB<br/>PostgreSQL)]
        PRODUCT_DB[(Product DB<br/>PostgreSQL)]
        REDIS[(Redis<br/>Cache & Basket)]
        PAYMENT_DB[(Payment DB<br/>PostgreSQL)]
    end
    
    WEB --> KRAKEND
    MOBILE --> KRAKEND
    
    KRAKEND --> USER
    KRAKEND --> PRODUCT
    KRAKEND --> BASKET
    KRAKEND --> PAYMENT
    
    PAYMENT --> KAFKA
    KAFKA --> BASKET
    KAFKA --> PRODUCT
    
    USER --> USER_DB
    PRODUCT --> PRODUCT_DB
    BASKET --> REDIS
    PAYMENT --> PAYMENT_DB
```

## Service Communication Flow

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
sequenceDiagram
    participant C as Client
    participant KG as KrakenD Gateway
    participant PS as Payment Service
    participant US as User Service
    participant PRS as Product Service
    participant BS as Basket Service
    participant K as Kafka
    
    C->>KG: Create Payment Request
    KG->>PS: POST /api/v1/payments
    PS->>US: Validate User (gRPC)
    US-->>PS: User Valid
    PS->>PRS: Validate Products (gRPC)
    PRS-->>PS: Products Valid
    PS->>BS: Validate Basket (gRPC)
    BS-->>PS: Basket Valid
    PS->>PS: Process Payment
    PS->>K: Publish PaymentCompleted Event
    K->>BS: Consume Event
    BS->>BS: Clear Basket
    K->>PRS: Consume Event
    PRS->>PRS: Update Stock
    PS-->>KG: Payment Response
    KG-->>C: Success Response
```

## Domain Layer Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "User Service Domain"
        U_ENTITY[User Entity]
        U_REPO[User Repository]
        U_SERVICE[User Service]
        U_ERR[Domain Errors]
    end
    
    subgraph "Product Service Domain"
        P_ENTITY[Product Entity]
        P_REPO[Product Repository]
        P_SERVICE[Product Service]
        P_ERR[Domain Errors]
    end
    
    subgraph "Basket Service Domain"
        B_ENTITY[Basket Entity]
        B_REPO[Basket Repository]
        B_SERVICE[Basket Service]
        B_ERR[Domain Errors]
    end
    
    subgraph "Payment Service Domain"
        PAY_ENTITY[Payment Entity]
        PAY_REPO[Payment Repository]
        PAY_GATEWAY[Payment Gateway]
        PAY_SERVICE[Payment Service]
        PAY_ERR[Domain Errors]
    end
    
    U_ENTITY --> U_REPO
    U_REPO --> U_SERVICE
    U_SERVICE --> U_ERR
    
    P_ENTITY --> P_REPO
    P_REPO --> P_SERVICE
    P_SERVICE --> P_ERR
    
    B_ENTITY --> B_REPO
    B_REPO --> B_SERVICE
    B_SERVICE --> B_ERR
    
    PAY_ENTITY --> PAY_REPO
    PAY_REPO --> PAY_SERVICE
    PAY_GATEWAY --> PAY_SERVICE
    PAY_SERVICE --> PAY_ERR
```

## CQRS Pattern Implementation

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph LR
    subgraph "Command Side"
        CMD[Command]
        CMD_HANDLER[Command Handler]
        CMD_REPO[Repository]
        CMD_DB[(Database)]
    end
    
    subgraph "Query Side"
        QUERY[Query]
        QUERY_HANDLER[Query Handler]
        QUERY_REPO[Repository]
        QUERY_DB[(Database)]
    end
    
    subgraph "Application Layer"
        APP_SERVICE[Application Service]
        DTO[DTOs]
    end
    
    CMD --> CMD_HANDLER
    CMD_HANDLER --> CMD_REPO
    CMD_REPO --> CMD_DB
    
    QUERY --> QUERY_HANDLER
    QUERY_HANDLER --> QUERY_REPO
    QUERY_REPO --> QUERY_DB
    
    APP_SERVICE --> CMD_HANDLER
    APP_SERVICE --> QUERY_HANDLER
    APP_SERVICE --> DTO
```

## Event-Driven Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Event Publishers"
        PS_PUB[Payment Service]
        US_PUB[User Service]
        PRS_PUB[Product Service]
        BS_PUB[Basket Service]
    end
    
    subgraph "Message Broker"
        KAFKA[Apache Kafka]
        TOPIC1[payment-events]
        TOPIC2[user-events]
        TOPIC3[product-events]
        TOPIC4[basket-events]
    end
    
    subgraph "Event Consumers"
        PS_CON[Payment Service]
        US_CON[User Service]
        PRS_CON[Product Service]
        BS_CON[Basket Service]
    end
    
    PS_PUB --> TOPIC1
    US_PUB --> TOPIC2
    PRS_PUB --> TOPIC3
    BS_PUB --> TOPIC4
    
    TOPIC1 --> PS_CON
    TOPIC2 --> US_CON
    TOPIC3 --> PRS_CON
    TOPIC4 --> BS_CON
    
    KAFKA --> TOPIC1
    KAFKA --> TOPIC2
    KAFKA --> TOPIC3
    KAFKA --> TOPIC4
```

## Payment Flow Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Payment Processing"
        CREATE[Create Payment]
        VALIDATE[Validate Payment]
        PROCESS[Process Payment]
        COMPLETE[Payment Complete]
    end
    
    subgraph "External Services"
        USER_SVC[User Service]
        PRODUCT_SVC[Product Service]
        BASKET_SVC[Basket Service]
        STRIPE[Stripe Gateway]
    end
    
    subgraph "Event Publishing"
        KAFKA_PUB[Kafka Publisher]
        PAYMENT_EVENT[PaymentCompleted Event]
        STOCK_EVENT[StockUpdated Event]
        BASKET_EVENT[BasketCleared Event]
    end
    
    CREATE --> VALIDATE
    VALIDATE --> USER_SVC
    VALIDATE --> PRODUCT_SVC
    VALIDATE --> BASKET_SVC
    VALIDATE --> PROCESS
    PROCESS --> STRIPE
    PROCESS --> COMPLETE
    COMPLETE --> KAFKA_PUB
    KAFKA_PUB --> PAYMENT_EVENT
    KAFKA_PUB --> STOCK_EVENT
    KAFKA_PUB --> BASKET_EVENT
```

## Technology Stack

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Backend Services"
        GO[Go 1.23]
        GIN[Gin Framework]
        GORM[GORM ORM]
        WIRE[Dependency Injection]
    end
    
    subgraph "Communication"
        GRPC[gRPC]
        HTTP[REST API]
        KAFKA[Apache Kafka]
        REDIS[Redis]
    end
    
    subgraph "Databases"
        POSTGRES[PostgreSQL]
        REDIS_DB[Redis Cache]
    end
    
    subgraph "Infrastructure"
        DOCKER[Docker]
        KRAKEND[KrakenD Gateway]
        STRIPE[Stripe Payment]
    end
    
    GO --> GIN
    GO --> GORM
    GO --> WIRE
    GO --> GRPC
    GO --> HTTP
    GO --> KAFKA
    
    GORM --> POSTGRES
    REDIS --> REDIS_DB
    
    DOCKER --> KRAKEND
    DOCKER --> STRIPE
```

## Dependency Management

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Dependabot Configuration"
        DEPENDABOT[Dependabot Bot]
        SCHEDULE[Weekly Schedule]
        AUTO_MERGE[Auto-merge PRs]
        LABELS[Auto-labeling]
    end
    
    subgraph "Dependency Ecosystems"
        GO_DEPS[Go Modules]
        NPM_DEPS[NPM Packages]
        DOCKER_DEPS[Docker Images]
        GITHUB_ACTIONS[GitHub Actions]
        TERRAFORM_DEPS[Terraform Modules]
        HELM_DEPS[Helm Charts]
    end
    
    subgraph "Update Process"
        CHECK[Check for Updates]
        CREATE_PR[Create Pull Request]
        RUN_TESTS[Run Tests]
        AUTO_MERGE_CHECK[Auto-merge Check]
        MERGE[Merge PR]
    end
    
    subgraph "Security & Quality"
        SECURITY_SCAN[Security Scanning]
        VULNERABILITY_CHECK[Vulnerability Check]
        COMPATIBILITY_TEST[Compatibility Test]
        BREAKING_CHANGE_CHECK[Breaking Change Check]
    end
    
    DEPENDABOT --> SCHEDULE
    DEPENDABOT --> AUTO_MERGE
    DEPENDABOT --> LABELS
    
    SCHEDULE --> CHECK
    CHECK --> CREATE_PR
    CREATE_PR --> RUN_TESTS
    RUN_TESTS --> AUTO_MERGE_CHECK
    AUTO_MERGE_CHECK --> MERGE
    
    GO_DEPS --> CHECK
    NPM_DEPS --> CHECK
    DOCKER_DEPS --> CHECK
    GITHUB_ACTIONS --> CHECK
    TERRAFORM_DEPS --> CHECK
    HELM_DEPS --> CHECK
    
    RUN_TESTS --> SECURITY_SCAN
    RUN_TESTS --> VULNERABILITY_CHECK
    RUN_TESTS --> COMPATIBILITY_TEST
    RUN_TESTS --> BREAKING_CHANGE_CHECK
```

## Service Dependencies

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TD
    subgraph "Core Services"
        USER[User Service]
        PRODUCT[Product Service]
        BASKET[Basket Service]
        PAYMENT[Payment Service]
    end
    
    subgraph "External Dependencies"
        STRIPE[Stripe API]
        KAFKA[Kafka Broker]
    end
    
    subgraph "Data Stores"
        USER_DB[(User DB)]
        PRODUCT_DB[(Product DB)]
        REDIS[(Redis)]
        PAYMENT_DB[(Payment DB)]
    end
    
    USER --> USER_DB
    PRODUCT --> PRODUCT_DB
    BASKET --> REDIS
    PAYMENT --> PAYMENT_DB
    
    PAYMENT --> USER
    PAYMENT --> PRODUCT
    PAYMENT --> BASKET
    PAYMENT --> STRIPE
    PAYMENT --> KAFKA
    
    KAFKA --> BASKET
    KAFKA --> PRODUCT
```

## API Gateway Routing

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Client Requests"
        WEB[Web Client]
        MOBILE[Mobile Client]
    end
    
    subgraph "KrakenD Gateway"
        GATEWAY[API Gateway<br/>Port: 8081]
    end
    
    subgraph "Service Endpoints"
        USER_EP[User Service<br/>Port: 8080]
        PRODUCT_EP[Product Service<br/>Port: 8081]
        BASKET_EP[Basket Service<br/>Port: 8083]
        PAYMENT_EP[Payment Service<br/>Port: 8084]
    end
    
    WEB --> GATEWAY
    MOBILE --> GATEWAY
    
    GATEWAY --> USER_EP
    GATEWAY --> PRODUCT_EP
    GATEWAY --> BASKET_EP
    GATEWAY --> PAYMENT_EP
```

## Database Schema Overview

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
erDiagram
    USER {
        uint id PK
        string email
        string password_hash
        string first_name
        string last_name
        string role
        timestamp created_at
        timestamp updated_at
    }
    
    PRODUCT {
        uint id PK
        string name
        text description
        decimal price
        int stock
        string category
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }
    
    BASKET {
        string id PK
        uint user_id FK
        timestamp created_at
        timestamp updated_at
        timestamp expires_at
    }
    
    BASKET_ITEM {
        string id PK
        string basket_id FK
        uint product_id FK
        int quantity
        decimal unit_price
        timestamp created_at
    }
    
    PAYMENT {
        string id PK
        uint user_id FK
        string order_id
        decimal amount
        string currency
        string status
        string payment_method
        string payment_provider
        string transaction_id
        timestamp created_at
        timestamp updated_at
        timestamp completed_at
        timestamp expires_at
    }
    
    PAYMENT_METHOD {
        string id PK
        uint user_id FK
        string type
        string provider
        string token
        string last_four_digits
        int expiry_month
        int expiry_year
        boolean is_default
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }
    
    USER ||--o{ BASKET : has
    BASKET ||--o{ BASKET_ITEM : contains
    PRODUCT ||--o{ BASKET_ITEM : referenced_by
    USER ||--o{ PAYMENT : makes
    USER ||--o{ PAYMENT_METHOD : owns
```

## Terraform Infrastructure Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "AWS Cloud Infrastructure"
        subgraph "VPC (10.0.0.0/16)"
            subgraph "Public Subnets"
                IGW[Internet Gateway]
                NAT1[NAT Gateway 1]
                NAT2[NAT Gateway 2]
                NAT3[NAT Gateway 3]
            end
            
            subgraph "Private Subnets"
                EKS[EKS Cluster]
                RDS1[PostgreSQL User]
                RDS2[PostgreSQL Product]
                RDS3[PostgreSQL Basket]
                RDS4[PostgreSQL Payment]
                REDIS[ElastiCache Redis]
                MSK[MSK Kafka Cluster]
            end
        end
        
        subgraph "Container Registry"
            ECR1[ECR User Service]
            ECR2[ECR Product Service]
            ECR3[ECR Basket Service]
            ECR4[ECR Payment Service]
        end
        
        subgraph "Load Balancer"
            ALB[Application Load Balancer]
            NLB[Network Load Balancer]
        end
    end
    
    subgraph "Kubernetes Cluster"
        subgraph "Control Plane"
            MASTER1[Master Node 1]
            MASTER2[Master Node 2]
            MASTER3[Master Node 3]
        end
        
        subgraph "Worker Nodes"
            WORKER1[Worker Node 1]
            WORKER2[Worker Node 2]
            WORKER3[Worker Node 3]
        end
        
        subgraph "Pods"
            USER_POD[User Service Pod]
            PRODUCT_POD[Product Service Pod]
            BASKET_POD[Basket Service Pod]
            PAYMENT_POD[Payment Service Pod]
            KRAKEND_POD[KrakenD Gateway Pod]
        end
    end
    
    IGW --> ALB
    ALB --> KRAKEND_POD
    KRAKEND_POD --> USER_POD
    KRAKEND_POD --> PRODUCT_POD
    KRAKEND_POD --> BASKET_POD
    KRAKEND_POD --> PAYMENT_POD
    
    USER_POD --> RDS1
    PRODUCT_POD --> RDS2
    BASKET_POD --> REDIS
    PAYMENT_POD --> RDS4
    PAYMENT_POD --> MSK
    
    ECR1 --> USER_POD
    ECR2 --> PRODUCT_POD
    ECR3 --> BASKET_POD
    ECR4 --> PAYMENT_POD
```

## Ansible Configuration Management

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Ansible Control Node"
        ANSIBLE[Ansible Controller]
        INVENTORY[Inventory Files]
        PLAYBOOKS[Playbooks]
        ROLES[Roles]
    end
    
    subgraph "Target Infrastructure"
        subgraph "Kubernetes Masters"
            K8S_MASTER1[Master Node 1]
            K8S_MASTER2[Master Node 2]
            K8S_MASTER3[Master Node 3]
        end
        
        subgraph "Kubernetes Workers"
            K8S_WORKER1[Worker Node 1]
            K8S_WORKER2[Worker Node 2]
            K8S_WORKER3[Worker Node 3]
        end
        
        subgraph "Database Servers"
            POSTGRES_DB[PostgreSQL Server]
            REDIS_DB[Redis Server]
            KAFKA_DB[Kafka Brokers]
        end
        
        subgraph "Monitoring Servers"
            PROMETHEUS_SRV[Prometheus Server]
            GRAFANA_SRV[Grafana Server]
            JAEGER_SRV[Jaeger Server]
        end
        
        subgraph "Load Balancers"
            NGINX_LB1[Nginx LB 1]
            NGINX_LB2[Nginx LB 2]
        end
    end
    
    subgraph "Ansible Roles"
        K8S_ROLE[k8s-setup Role]
        HELM_ROLE[helm-deploy Role]
        MONITORING_ROLE[monitoring Role]
        SECURITY_ROLE[security Role]
        BACKUP_ROLE[backup Role]
    end
    
    ANSIBLE --> INVENTORY
    ANSIBLE --> PLAYBOOKS
    ANSIBLE --> ROLES
    
    PLAYBOOKS --> K8S_ROLE
    PLAYBOOKS --> HELM_ROLE
    PLAYBOOKS --> MONITORING_ROLE
    PLAYBOOKS --> SECURITY_ROLE
    PLAYBOOKS --> BACKUP_ROLE
    
    K8S_ROLE --> K8S_MASTER1
    K8S_ROLE --> K8S_MASTER2
    K8S_ROLE --> K8S_MASTER3
    K8S_ROLE --> K8S_WORKER1
    K8S_ROLE --> K8S_WORKER2
    K8S_ROLE --> K8S_WORKER3
    
    HELM_ROLE --> K8S_MASTER1
    
    MONITORING_ROLE --> PROMETHEUS_SRV
    MONITORING_ROLE --> GRAFANA_SRV
    MONITORING_ROLE --> JAEGER_SRV
    
    SECURITY_ROLE --> K8S_MASTER1
    SECURITY_ROLE --> K8S_MASTER2
    SECURITY_ROLE --> K8S_MASTER3
    SECURITY_ROLE --> K8S_WORKER1
    SECURITY_ROLE --> K8S_WORKER2
    SECURITY_ROLE --> K8S_WORKER3
    
    BACKUP_ROLE --> K8S_MASTER1
    BACKUP_ROLE --> POSTGRES_DB
    BACKUP_ROLE --> REDIS_DB
```

## Kubernetes Cluster Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Kubernetes Control Plane"
        subgraph "Master Node 1"
            API1[API Server]
            ETCD1[etcd]
            CM1[Controller Manager]
            SCHED1[Scheduler]
        end
        
        subgraph "Master Node 2"
            API2[API Server]
            ETCD2[etcd]
            CM2[Controller Manager]
            SCHED2[Scheduler]
        end
        
        subgraph "Master Node 3"
            API3[API Server]
            ETCD3[etcd]
            CM3[Controller Manager]
            SCHED3[Scheduler]
        end
    end
    
    subgraph "Worker Nodes"
        subgraph "Worker Node 1"
            KUBELET1[kubelet]
            PROXY1[kube-proxy]
            CONTAINERD1[containerd]
            CNI1[CNI Plugin]
        end
        
        subgraph "Worker Node 2"
            KUBELET2[kubelet]
            PROXY2[kube-proxy]
            CONTAINERD2[containerd]
            CNI2[CNI Plugin]
        end
        
        subgraph "Worker Node 3"
            KUBELET3[kubelet]
            PROXY3[kube-proxy]
            CONTAINERD3[containerd]
            CNI3[CNI Plugin]
        end
    end
    
    subgraph "Pods and Services"
        subgraph "ddd-micro Namespace"
            USER_SVC[User Service]
            PRODUCT_SVC[Product Service]
            BASKET_SVC[Basket Service]
            PAYMENT_SVC[Payment Service]
            KRAKEND_SVC[KrakenD Gateway]
        end
        
        subgraph "System Pods"
            DNS[CoreDNS]
            CNI_POD[CNI Pods]
            PROXY_POD[Proxy Pods]
        end
    end
    
    subgraph "Storage"
        PV1[Persistent Volume 1]
        PV2[Persistent Volume 2]
        PV3[Persistent Volume 3]
        SC[Storage Classes]
    end
    
    subgraph "Network"
        CNI_NET[CNI Network]
        SERVICE_NET[Service Network]
        POD_NET[Pod Network]
    end
    
    API1 --> KUBELET1
    API1 --> KUBELET2
    API1 --> KUBELET3
    
    API2 --> KUBELET1
    API2 --> KUBELET2
    API2 --> KUBELET3
    
    API3 --> KUBELET1
    API3 --> KUBELET2
    API3 --> KUBELET3
    
    KUBELET1 --> CONTAINERD1
    KUBELET2 --> CONTAINERD2
    KUBELET3 --> CONTAINERD3
    
    CONTAINERD1 --> USER_SVC
    CONTAINERD1 --> PRODUCT_SVC
    CONTAINERD2 --> BASKET_SVC
    CONTAINERD2 --> PAYMENT_SVC
    CONTAINERD3 --> KRAKEND_SVC
    
    CNI1 --> CNI_NET
    CNI2 --> CNI_NET
    CNI3 --> CNI_NET
    
    USER_SVC --> PV1
    PRODUCT_SVC --> PV2
    BASKET_SVC --> PV3
```

## CI/CD Pipeline Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Source Control"
        GITHUB[GitHub Repository]
        MAIN[main branch]
        DEVELOP[develop branch]
        FEATURE[feature branches]
    end
    
    subgraph "CI/CD Triggers"
        PUSH[Push Events]
        PR[Pull Request Events]
        SCHEDULE[Scheduled Events]
    end
    
    subgraph "CI Pipeline"
        CHECKOUT[Checkout Code]
        TEST[Run Tests]
        BUILD[Build Images]
        SECURITY[Security Scan]
        LINT[Code Linting]
    end
    
    subgraph "CD Pipeline"
        subgraph "Staging Environment"
            TF_STAGING[Terraform Apply - Staging]
            ANSIBLE_STAGING[Ansible Deploy - Staging]
            HELM_STAGING[Helm Deploy - Staging]
        end
        
        subgraph "Production Environment"
            TF_PROD[Terraform Apply - Prod]
            ANSIBLE_PROD[Ansible Deploy - Prod]
            HELM_PROD[Helm Deploy - Prod]
        end
    end
    
    subgraph "Infrastructure as Code"
        TERRAFORM[Terraform Modules]
        ANSIBLE[Ansible Playbooks]
        HELM[Helm Charts]
    end
    
    subgraph "Container Registry"
        ECR[Amazon ECR]
        DOCKERHUB[Docker Hub]
    end
    
    subgraph "Deployment Targets"
        K8S_STAGING[K8s Staging Cluster]
        K8S_PROD[K8s Production Cluster]
    end
    
    subgraph "Monitoring & Notifications"
        SLACK[Slack Notifications]
        EMAIL[Email Alerts]
        DASHBOARD[Monitoring Dashboard]
    end
    
    GITHUB --> PUSH
    GITHUB --> PR
    GITHUB --> SCHEDULE
    
    PUSH --> CHECKOUT
    PR --> CHECKOUT
    SCHEDULE --> CHECKOUT
    
    CHECKOUT --> TEST
    CHECKOUT --> LINT
    TEST --> BUILD
    LINT --> BUILD
    BUILD --> SECURITY
    
    SECURITY --> TF_STAGING
    SECURITY --> TF_PROD
    
    TF_STAGING --> ANSIBLE_STAGING
    ANSIBLE_STAGING --> HELM_STAGING
    HELM_STAGING --> K8S_STAGING
    
    TF_PROD --> ANSIBLE_PROD
    ANSIBLE_PROD --> HELM_PROD
    HELM_PROD --> K8S_PROD
    
    BUILD --> ECR
    BUILD --> DOCKERHUB
    
    K8S_STAGING --> SLACK
    K8S_PROD --> SLACK
    K8S_STAGING --> EMAIL
    K8S_PROD --> EMAIL
    
    K8S_STAGING --> DASHBOARD
    K8S_PROD --> DASHBOARD
```

## Deployment Architecture

```mermaid
%%{init: {'theme':'dark', 'themeVariables': { 'primaryColor': '#ff6b6b', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#ff6b6b', 'lineColor': '#ffffff', 'secondaryColor': '#4ecdc4', 'tertiaryColor': '#45b7d1', 'background': '#2c3e50', 'mainBkg': '#34495e', 'secondBkg': '#2c3e50', 'tertiaryBkg': '#34495e'}}}%%
graph TB
    subgraph "Docker Containers"
        USER_CONTAINER[User Service Container]
        PRODUCT_CONTAINER[Product Service Container]
        BASKET_CONTAINER[Basket Service Container]
        PAYMENT_CONTAINER[Payment Service Container]
        GATEWAY_CONTAINER[KrakenD Gateway Container]
        KAFKA_CONTAINER[Kafka Container]
        ZOOKEEPER_CONTAINER[Zookeeper Container]
        REDIS_CONTAINER[Redis Container]
        POSTGRES_CONTAINER[PostgreSQL Container]
    end
    
    subgraph "Docker Network"
        NETWORK[ddd-micro-network]
    end
    
    USER_CONTAINER --> NETWORK
    PRODUCT_CONTAINER --> NETWORK
    BASKET_CONTAINER --> NETWORK
    PAYMENT_CONTAINER --> NETWORK
    GATEWAY_CONTAINER --> NETWORK
    KAFKA_CONTAINER --> NETWORK
    ZOOKEEPER_CONTAINER --> NETWORK
    REDIS_CONTAINER --> NETWORK
    POSTGRES_CONTAINER --> NETWORK
```