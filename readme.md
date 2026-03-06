# Rental Car API

A robust, secure, and scalable REST API for a Car Rental Service built with Go (Golang). This project provides a complete solution for managing car rentals, user deposits, automated background tasks, and administrative insights.

## 🚀 Tech Stack

- **Language:** Go (1.21+)
- **Web Framework:** [Echo](https://echo.labstack.com/)
- **Database:** PostgreSQL
- **Database Driver:** [pgx/v5](https://github.com/jackc/pgx)
- **Authentication:** JWT (JSON Web Tokens)
- **API Documentation:** [Swaggo](https://github.com/swaggo/swag) (Swagger 2.0)
- **Image Hosting:** [Cloudinary](https://cloudinary.com/)
- **Payment Gateway:** [Xendit](https://www.xendit.co/)
- **Email Service:** Gmail API (via Google OAuth2)
- **Testing:** [testify](https://github.com/stretchr/testify) & Mocking

## ✨ Core Features

### 🔐 Authentication & User Management
- Secure registration and login with bcrypt password hashing.
- JWT-based authentication with protected routes.
- Profile management and password reset functionality.
- Role-based Access Control (RBAC): `user` and `admin` roles.

### 🚗 Car Management
- Full CRUD operations for cars (Admin only).
- **Advanced Search:** Full-text search by name or description.
- **Filtering:** Filter by category and availability.
- **Pagination & Sorting:** Efficient data retrieval with limit/offset and multi-field sorting.
- **Image Upload:** Direct file upload to Cloudinary for car images.

### 💳 Rental & Payment System
- **Atomic Transactions:** Ensures data integrity during the rental process (car stock vs. user deposit).
- **Payment Gateway:** Integration with Xendit for automated invoicing.
- **Webhook Security:** Secure callback verification using Xendit callback tokens.
- Internal deposit system for seamless internal payments.

### 🤖 Automated Tasks (Background Worker)
- **Rental Expiration Worker:** An hourly background job that:
    - Automatically marks rentals as `overdue` when past their due date.
    - Sends automated email reminders to users with overdue rentals.

### 📊 Admin Dashboard
- Real-time statistics including Total Revenue, Total Rentals, and Total Users.
- Top 5 most popular cars based on rental history.

## 📂 Project Structure

```text
.
├── app/                # Application entry point & configuration
│   ├── config/         # Environment configuration loader
│   ├── dto/            # Data Transfer Objects
│   ├── handler/        # HTTP Handlers (Controller layer)
│   └── middleware/     # Custom middlewares (JWT, Validator)
├── database/           # Database connection & adapters
├── docs/               # Swagger generated documentation
├── model/              # Database models (Domain entities)
├── repository/         # Data Access Layer (SQL queries)
├── service/            # Business Logic Layer
└── main.go             # Main application file
```

## 🛠️ Getting Started

### Prerequisites
- Go 1.21 or later
- PostgreSQL database
- Cloudinary account (for image uploads)
- Xendit account (for payments)
- Gmail API credentials (for email notifications)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd car-rental-miniproject
   ```

2. **Setup environment variables:**
   Copy `.env.example` to `.env` and fill in your credentials:
   ```bash
   cp .env.example .env
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the database schema:**
   Execute the queries in `ddl.sql` on your PostgreSQL instance.

5. **Start the application:**
   ```bash
   go run main.go
   ```

## 📖 API Documentation

Once the server is running, you can access the interactive Swagger documentation at:
`http://localhost:8080/swagger/index.html`

### Key Endpoints Summary

| Method | Endpoint | Description | Access |
| :--- | :--- | :--- | :--- |
| POST | `/api/auth/register` | Register a new account | Public |
| POST | `/api/auth/login` | Login and get JWT token | Public |
| GET | `/api/cars` | Get cars with search & pagination | Public |
| POST | `/api/cars/upload` | Upload car image | Admin |
| POST | `/api/rentals` | Rent a car | User |
| GET | `/api/admin/dashboard` | Get dashboard statistics | Admin |

## 🧪 Testing

The project uses a mocking strategy to test services in isolation. To run the tests:

```bash
go test ./service/...
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
