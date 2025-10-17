# DDD Microservices Architecture

## System Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[Web Client]
        MOBILE[Mobile Client]
        API_CLIENT[API Client]
    end
    
    subgraph "API Gateway"
        KRAKEND[KrakenD Gateway]
    end
    
    subgraph "Microservices"
        USER_SERVICE[User Service<br/>Port: 8080<br/>gRPC: 9090]
        PRODUCT_SERVICE[Product Service<br/>Port: 8081<br/>gRPC: 9091]
        BASKET_SERVICE[Basket Service<br/>Port: 8083<br/>Redis]
        PAYMENT_SERVICE[Payment Service<br/>Port: 8084]
    end
    
    subgraph "Databases"
        POSTGRES_USER[(PostgreSQL<br/>User DB)]
        POSTGRES_PRODUCT[(PostgreSQL<br/>Product DB)]
        REDIS[(Redis<br/>Basket Cache)]
        POSTGRES_PAYMENT[(PostgreSQL<br/>Payment DB)]
    end
    
    WEB --> KRAKEND
    MOBILE --> KRAKEND
    API_CLIENT --> KRAKEND
    
    KRAKEND --> USER_SERVICE
    KRAKEND --> PRODUCT_SERVICE
    KRAKEND --> BASKET_SERVICE
    KRAKEND --> PAYMENT_SERVICE
    
    USER_SERVICE --> POSTGRES_USER
    PRODUCT_SERVICE --> POSTGRES_PRODUCT
    BASKET_SERVICE --> REDIS
    PAYMENT_SERVICE --> POSTGRES_PAYMENT
    
    BASKET_SERVICE -.->|gRPC| USER_SERVICE
    BASKET_SERVICE -.->|gRPC| PRODUCT_SERVICE
    PRODUCT_SERVICE -.->|gRPC| USER_SERVICE
```

## Service Communication

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant BasketService
    participant UserService
    participant ProductService
    participant Redis
    
    Client->>Gateway: POST /api/v1/users/basket/items
    Gateway->>BasketService: Forward Request
    
    BasketService->>UserService: Validate User Token
    UserService-->>BasketService: User Info
    
    BasketService->>ProductService: Get Product Details
    ProductService-->>BasketService: Product Info & Stock
    
    BasketService->>Redis: Store/Update Basket
    Redis-->>BasketService: Success
    
    BasketService-->>Gateway: Response
    Gateway-->>Client: Response
```

## Domain Architecture

```mermaid
graph TD
    subgraph "User Service"
        USER_DOMAIN[User Domain]
        USER_APP[User Application]
        USER_INFRA[User Infrastructure]
        USER_INTERFACE[User Interface]
    end
    
    subgraph "Product Service"
        PRODUCT_DOMAIN[Product Domain]
        PRODUCT_APP[Product Application]
        PRODUCT_INFRA[Product Infrastructure]
        PRODUCT_INTERFACE[Product Interface]
    end
    
    subgraph "Basket Service"
        BASKET_DOMAIN[Basket Domain]
        BASKET_APP[Basket Application]
        BASKET_INFRA[Basket Infrastructure]
        BASKET_INTERFACE[Basket Interface]
    end
    
    USER_DOMAIN --> USER_APP
    USER_APP --> USER_INFRA
    USER_APP --> USER_INTERFACE
    
    PRODUCT_DOMAIN --> PRODUCT_APP
    PRODUCT_APP --> PRODUCT_INFRA
    PRODUCT_APP --> PRODUCT_INTERFACE
    
    BASKET_DOMAIN --> BASKET_APP
    BASKET_APP --> BASKET_INFRA
    BASKET_APP --> BASKET_INTERFACE
```

## CQRS Pattern

```mermaid
graph LR
    subgraph "Command Side"
        CMD[Command]
        CMD_HANDLER[Command Handler]
        DOMAIN[Domain Logic]
        REPO[Repository]
    end
    
    subgraph "Query Side"
        QUERY[Query]
        QUERY_HANDLER[Query Handler]
        READ_MODEL[Read Model]
    end
    
    CMD --> CMD_HANDLER
    CMD_HANDLER --> DOMAIN
    DOMAIN --> REPO
    
    QUERY --> QUERY_HANDLER
    QUERY_HANDLER --> READ_MODEL
```

## Technology Stack

```mermaid
graph TB
    subgraph "Backend"
        GO[Go 1.23]
        GIN[Gin Framework]
        GRPC[gRPC]
        WIRE[Wire DI]
    end
    
    subgraph "Databases"
        POSTGRES[PostgreSQL]
        REDIS_CACHE[Redis]
    end
    
    subgraph "Documentation"
        SWAGGER[Swagger/OpenAPI]
        PROTO[Protocol Buffers]
    end
    
    subgraph "Infrastructure"
        DOCKER[Docker]
        COMPOSE[Docker Compose]
        KRAKEND[KrakenD]
    end
    
    GO --> GIN
    GO --> GRPC
    GO --> WIRE
    GIN --> SWAGGER
    GRPC --> PROTO
    GO --> POSTGRES
    GO --> REDIS_CACHE
    GO --> DOCKER
    GO --> COMPOSE
    GO --> KRAKEND
```

## API Endpoints

```mermaid
graph TB
    subgraph "User Service"
        USER_AUTH[Authentication]
        USER_PROFILE[User Profile]
        USER_ADMIN[User Management]
    end
    
    subgraph "Product Service"
        PRODUCT_PUBLIC[Public Products]
        PRODUCT_USER[User Products]
        PRODUCT_ADMIN[Admin Products]
    end
    
    subgraph "Basket Service"
        BASKET_USER[User Basket]
        BASKET_ADMIN[Admin Basket]
    end
    
    USER_AUTH --> USER_PROFILE
    USER_PROFILE --> USER_ADMIN
    PRODUCT_PUBLIC --> PRODUCT_USER
    PRODUCT_USER --> PRODUCT_ADMIN
    BASKET_USER --> BASKET_ADMIN
```

## Security & Authentication

```mermaid
graph TB
    subgraph "Authentication Flow"
        LOGIN[User Login]
        JWT[JWT Token]
        VALIDATE[Token Validation]
        RBAC[Role-Based Access]
    end
    
    subgraph "Service Communication"
        INTERCEPTOR[gRPC Interceptor]
        MIDDLEWARE[HTTP Middleware]
        AUTH_CHECK[Auth Check]
    end
    
    LOGIN --> JWT
    JWT --> VALIDATE
    VALIDATE --> RBAC
    RBAC --> INTERCEPTOR
    RBAC --> MIDDLEWARE
    INTERCEPTOR --> AUTH_CHECK
    MIDDLEWARE --> AUTH_CHECK
```

## Deployment Architecture

```mermaid
graph TB
    subgraph "Development"
        DEV_ENV[Local Development]
        DEV_DB[Local Databases]
    end
    
    subgraph "Production"
        PROD_ENV[Production Environment]
        PROD_DB[Production Databases]
        LOAD_BALANCER[Load Balancer]
    end
    
    DEV_ENV --> DEV_DB
    PROD_ENV --> PROD_DB
    PROD_ENV --> LOAD_BALANCER
```

## License

MIT License

Copyright (c) 2024 DDD Microservices Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
