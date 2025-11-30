# Gram Panchayat Management System

A modern, full-stack web application for digital governance of Gram Panchayat operations.

## 🚀 Features

### For Citizens
- Online application for certificates (Birth, Death, Income, Caste, Residence)
- Complaint registration and tracking
- Property tax payment
- View notices and announcements
- Track application status
- Download certificates

### For Admins
- User management
- Application approval workflow
- Complaint management
- Property tax management
- Notice and meeting management
- Financial reports and analytics
- Dashboard with real-time statistics

## 🛠️ Tech Stack

### Backend
- **Language:** Go (Golang) 1.21+
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL 15
- **Authentication:** JWT
- **ORM:** GORM

### Frontend
- **Framework:** React 18+
- **Build Tool:** Vite
- **Styling:** Tailwind CSS
- **Icons:** Lucide React
- **HTTP Client:** Axios

## 📋 Prerequisites

Before you begin, ensure you have installed:
- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL 15+
- Docker & Docker Compose (optional)

## 🔧 Installation

### Option 1: Manual Setup

#### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/gram-panchayat.git
cd gram-panchayat/backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create `.env` file:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Create PostgreSQL database:
```bash
createdb gram_panchayat
```

5. Run the application:
```bash
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

#### Frontend Setup

1. Navigate to frontend directory:
```bash
cd ../frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create `.env` file:
```bash
cp .env.example .env
# Update VITE_API_URL if needed
```

4. Start development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

### Option 2: Docker Setup

1. Make sure Docker and Docker Compose are installed

2. Clone the repository:
```bash
git clone https://github.com/yourusername/gram-panchayat.git
cd gram-panchayat
```

3. Start all services:
```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Backend API on port 8080
- Frontend on port 5173
- PgAdmin on port 5050 (optional)

4. Access the application:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- PgAdmin: http://localhost:5050

## 📁 Project Structure

```
gram-panchayat/
├── backend/                 # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go     # Entry point
│   ├── internal/
│   │   ├── config/         # Configuration
│   │   ├── database/       # Database connection
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Middleware
│   │   ├── models/         # Data models
│   │   ├── repository/     # Database operations
│   │   ├── service/        # Business logic
│   │   └── utils/          # Utility functions
│   ├── uploads/            # File uploads
│   ├── .env.example
│   ├── go.mod
│   └── Dockerfile
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   ├── hooks/         # Custom hooks
│   │   └── utils/         # Utility functions
│   ├── package.json
│   └── Dockerfile
└── docker-compose.yml
```

## 🔐 Default Credentials

After first setup, use these credentials to login as admin:

- **Email:** admin@grampanchayat.gov.in
- **Password:** admin@123

**⚠️ Important:** Change the default password immediately after first login!

## 🌐 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new citizen
- `POST /api/auth/login` - Login
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password
- `GET /api/auth/profile` - Get user profile

### Applications
- `POST /api/services/apply` - Submit application
- `GET /api/services/applications` - List applications
- `GET /api/services/applications/:id` - Get application details
- `PUT /api/admin/applications/:id/status` - Update status (admin)

### Complaints
- `POST /api/complaints` - Create complaint
- `GET /api/complaints` - List complaints
- `GET /api/complaints/:id` - Get complaint details
- `PUT /api/complaints/:id` - Update complaint (admin)

### Property Tax
- `GET /api/property-tax/properties` - List properties
- `POST /api/property-tax/:propertyId/payment` - Make payment
- `GET /api/property-tax/payment-history` - Payment history

### Admin
- `GET /api/admin/users` - List all users
- `GET /api/dashboard/admin` - Admin dashboard stats

Full API documentation available at: `/docs/API.md`

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 🚀 Deployment

### Backend Deployment

1. Build the binary:
```bash
cd backend
go build -o main ./cmd/server
```

2. Set environment variables on your server

3. Run the binary:
```bash
./main
```

### Frontend Deployment

1. Build for production:
```bash
cd frontend
npm run build
```

2. Deploy the `dist` folder to your hosting service (Netlify, Vercel, etc.)

### Docker Deployment

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 📝 Environment Variables

### Backend (.env)
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=gram_panchayat
JWT_SECRET=your-secret-key
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080/api
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👥 Support

For support, email support@grampanchayat.gov.in or create an issue in the GitHub repository.

## 🙏 Acknowledgments

- Gin Web Framework
- React Team
- Tailwind CSS
- GORM
- All contributors

---

Made with ❤️ for Digital India Initiative