# SpotSync API

Smart Parking and EV Charging Reservation System

## Live URL

https://go-spotsync-api.onrender.com

## Tech Stack

- Go 1.26
- Echo v4 (Web Framework)
- GORM (ORM)
- PostgreSQL / NeonDB (Database)
- JWT (Authentication)
- bcrypt (Password Hashing)

## Architecture

This project follows Clean Architecture. Each domain is separated into 4 layers.

Handler receives the HTTP request and calls the Service.
Service contains the business logic and calls the Repository.
Repository handles all database operations using GORM.
DTO defines the request and response structures.

Dependency injection is done manually in main.go.

## How to Run Locally

Step 1: Clone the repository

git clone https://github.com/shuvoredsky/go-spotsync
cd go-spotsync

Step 2: Create a .env file in the project root

PORT=8080
DSN=postgresql://username:password@host/dbname?sslmode=require
JWT_SECRET=your_secret_key

Step 3: Run the server

go run cmd/main.go

## Environment Variables

PORT - The port the server runs on
DSN - PostgreSQL connection string
JWT_SECRET - Secret key for signing JWT tokens

## API Endpoints

### Authentication

POST /api/v1/auth/register - Register a new user (Public)
POST /api/v1/auth/login - Login and get JWT token (Public)

### Parking Zones

GET /api/v1/zones - Get all parking zones (Public)
GET /api/v1/zones/:id - Get a single parking zone (Public)
POST /api/v1/zones - Create a parking zone (Admin only)
PUT /api/v1/zones/:id - Update a parking zone (Admin only)
DELETE /api/v1/zones/:id - Delete a parking zone (Admin only)

### Reservations

POST /api/v1/reservations - Create a reservation (Authenticated)
GET /api/v1/reservations/my-reservations - Get my reservations (Authenticated)
DELETE /api/v1/reservations/:id - Cancel a reservation (Authenticated)
GET /api/v1/reservations - Get all reservations (Admin only)

## Authentication

After login, you will receive a JWT token.
Add it to the Authorization header like this:

Authorization: Bearer your_token_here

## Concurrency Handling

When two users try to book the last available EV spot at the same time, a race condition can occur.
This project solves that problem using PostgreSQL transactions combined with row-level locking (SELECT FOR UPDATE).
This ensures only one reservation is created even under high concurrency.

## User Roles

driver - Can register, login, view zones, create and cancel their own reservations
admin - Has all driver permissions plus can manage parking zones and view all reservations
