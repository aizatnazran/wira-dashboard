# WIRA Dashboard

The WIRA Dashboard is a web application designed to manage and display comprehensive score data for players of the game WIRA. This application combines an efficient backend capable of handling over 100,000 lines of data with an interactive, user-friendly frontend that delivers a seamless experience for end users.

## Features

### Core Features
- Interactive dashboard displaying player rankings and scores
- Support for 8 different character classes
- Search functionality
- Efficient pagination handling 100,000+ records
- Responsive design for both desktop and mobile
- Data caching using Redis for optimal performance

### Optional Features
- Secure login system
- Two-factor authentication (2FA)
- Session management
- Password encryption

## Tech Stack

- Frontend: Vue.js
- Backend: Golang
- Database: PostgreSQL
- Server: Nginx

## Prerequisites

Before you begin, ensure you have the following installed:
- Node.js (v16 or higher)
- Go (v1.21 or higher)
- PostgreSQL (v15 or higher)
- Git

## Getting Started

### 1. Clone the Repository
```bash
git clone https://github.com/aizatnazran/wira-dashboard.git
cd wira-dashboard
```

### 2. Database Setup
Create a PostgreSQL database:
```bash
createdb wira
```
Navigate to the database migration directory:
```bash
cd backend/db
```
Run the database migrations and seed data:
```bash
go run migrate.go --reset
```

### 3. Backend Setup
Navigate to the backend directory:
```bash
cd backend
```
Create a .env file with the required environment variables:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=wira
DB_SSLMODE=disable
JWT_SECRET=your_jwt_secret
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

Install Go dependencies:
```bash
go mod download
go mod tidy
```
Start the backend server:
```bash
go run main.go
```

### 4. Frontend Setup
Navigate to the frontend directory:
```bash
cd frontend
```
Install dependencies:
```bash
npm install
```
Start the development server:
```bash
npm run serve
```


### 5. Project Structure
├── frontend/               # Vue.js frontend application
│   ├── src/                # Source files
│   │   ├── assets/         # Static assets
│   │   ├── components/     # Vue components
│   │   ├── views/          # Page components
│   │   └── api/            # API utilities
│   └── package.json        # Frontend dependencies
│
├── backend/               # Golang backend server
│   ├── auth/              # Authentication package
│   ├── cache/             # Redis caching package
│   ├── config/            # Configuration package
│   ├── db/                # Database migrations
│   ├── ranking/           # Ranking logic package
│   └── main.go            # Entry point


## API Endpoints
The backend provides RESTful APIs with the following main endpoints:

POST /api/auth/register - Register a new user
POST /api/auth/login - Login user
GET /api/profile - Get user profile
GET /api/rankings - Get rankings with pagination
GET /api/classes - Get available classes