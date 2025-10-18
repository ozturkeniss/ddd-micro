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