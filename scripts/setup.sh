#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}   Gram Panchayat System Setup Script${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""

# Check if Go is installed
echo -e "${YELLOW}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go $(go version) found${NC}"

# Check if Node.js is installed
echo -e "${YELLOW}Checking Node.js installation...${NC}"
if ! command -v node &> /dev/null; then
    echo -e "${RED}Node.js is not installed. Please install Node.js 18+ first.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Node.js $(node --version) found${NC}"

# Check if PostgreSQL is installed
echo -e "${YELLOW}Checking PostgreSQL installation...${NC}"
if ! command -v psql &> /dev/null; then
    echo -e "${RED}PostgreSQL is not installed. Please install PostgreSQL 15+ first.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ PostgreSQL found${NC}"

echo ""
echo -e "${YELLOW}Setting up Backend...${NC}"

# Navigate to backend directory
cd backend

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file from example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ .env file created${NC}"
    echo -e "${YELLOW}Please update the .env file with your configuration${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Download Go dependencies
echo -e "${YELLOW}Downloading Go dependencies...${NC}"
go mod download
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Go dependencies installed${NC}"
else
    echo -e "${RED}✗ Failed to install Go dependencies${NC}"
    exit 1
fi

# Create uploads directories
echo -e "${YELLOW}Creating uploads directories...${NC}"
mkdir -p uploads/documents uploads/certificates uploads/profiles
echo -e "${GREEN}✓ Upload directories created${NC}"

echo ""
echo -e "${YELLOW}Setting up Frontend...${NC}"

# Navigate to frontend directory
cd ../frontend

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating frontend .env file...${NC}"
    echo "VITE_API_URL=http://localhost:8080/api" > .env
    echo -e "${GREEN}✓ Frontend .env file created${NC}"
else
    echo -e "${GREEN}✓ Frontend .env file already exists${NC}"
fi

# Install npm dependencies
echo -e "${YELLOW}Installing npm dependencies...${NC}"
npm install
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ npm dependencies installed${NC}"
else
    echo -e "${RED}✗ Failed to install npm dependencies${NC}"
    exit 1
fi

cd ..

echo ""
echo -e "${YELLOW}Setting up Database...${NC}"

# Prompt for database configuration
read -p "Enter PostgreSQL username [postgres]: " DB_USER
DB_USER=${DB_USER:-postgres}

read -sp "Enter PostgreSQL password: " DB_PASSWORD
echo ""

read -p "Enter database name [gram_panchayat]: " DB_NAME
DB_NAME=${DB_NAME:-gram_panchayat}

# Create database
echo -e "${YELLOW}Creating database...${NC}"
PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -h localhost -c "CREATE DATABASE $DB_NAME;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database created successfully${NC}"
else
    echo -e "${YELLOW}⚠ Database might already exist or creation failed${NC}"
fi

echo ""
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}   Setup Complete!${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
echo "1. Update backend/.env with your configuration"
echo "2. Start the backend server:"
echo -e "   ${GREEN}cd backend && go run cmd/server/main.go${NC}"
echo ""
echo "3. In a new terminal, start the frontend:"
echo -e "   ${GREEN}cd frontend && npm run dev${NC}"
echo ""
echo "4. Access the application:"
echo -e "   Frontend: ${GREEN}http://localhost:5173${NC}"
echo -e "   Backend API: ${GREEN}http://localhost:8080${NC}"
echo ""
echo "5. Default admin credentials:"
echo -e "   Email: ${GREEN}admin@grampanchayat.gov.in${NC}"
echo -e "   Password: ${GREEN}admin@123${NC}"
echo -e "   ${RED}⚠ Change this password immediately after first login!${NC}"
echo ""
echo -e "${YELLOW}Or use Docker:${NC}"
echo -e "   ${GREEN}docker-compose up -d${NC}"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"