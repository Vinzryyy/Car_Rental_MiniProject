# Database Entity Relationship Diagram (ERD)

```mermaid
erDiagram
    USERS ||--o{ RENTAL_HISTORIES : makes
    USERS ||--o{ TOP_UP_TRANSACTIONS : performs
    USERS ||--o{ USER_SESSIONS : has
    CARS ||--o{ RENTAL_HISTORIES : "rented in"

    USERS {
        uuid id PK
        string email UK
        string password
        decimal deposit_amount
        string role
        timestamp created_at
        timestamp updated_at
    }

    CARS {
        uuid id PK
        string name
        boolean availability
        int stock_availability
        decimal rental_costs
        string category
        text description
        string image_url
        timestamp created_at
        timestamp updated_at
    }

    RENTAL_HISTORIES {
        uuid id PK
        uuid user_id FK
        uuid car_id FK
        timestamp rental_date
        timestamp return_date
        decimal total_cost
        string status
        string payment_status
        string payment_url
        timestamp created_at
        timestamp updated_at
    }

    TOP_UP_TRANSACTIONS {
        uuid id PK
        uuid user_id FK
        decimal amount
        string status
        string payment_method
        string payment_url
        timestamp created_at
        timestamp updated_at
    }

    USER_SESSIONS {
        uuid id PK
        uuid user_id FK
        string token UK
        timestamp expires_at
        timestamp created_at
    }
```

## Description of Entities

### Users
Stores user account information, including authentication credentials, current deposit balance, and roles (user/admin).

### Cars
Main entity representing the rental products. Tracks stock levels, pricing, and categories.

### Rental Histories
Tracks all car rental transactions, including the duration, cost, and payment status. Links users to the cars they rent.

### Top-Up Transactions
Records all deposit addition attempts, whether pending or completed via the payment gateway.

### User Sessions
Manages active JWT sessions and refresh tokens for secure authentication.
