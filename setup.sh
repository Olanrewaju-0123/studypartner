#!/bin/bash

echo "🚀 Setting up AI Study Partner..."

# Check if required tools are installed
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    else
        echo "✅ $1 is installed"
    fi
}

echo "📋 Checking prerequisites..."
check_command "go"
check_command "node"
check_command "npm"

# Setup backend
echo "🔧 Setting up backend..."
cd backend

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    cp env.example .env
    echo "📝 Created .env file. Please edit it with your configuration."
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

cd ..

# Setup frontend
echo "🎨 Setting up frontend..."
cd frontend

# Copy environment file if it doesn't exist
if [ ! -f .env.local ]; then
    cp env.example .env.local
    echo "📝 Created .env.local file. Please edit it with your configuration."
fi

# Install Node.js dependencies
echo "📦 Installing Node.js dependencies..."
npm install

cd ..

echo "✅ Setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Set up your database (see SETUP.md for details)"
echo "2. Configure your .env files"
echo "3. Start the backend: cd backend && go run cmd/main.go"
echo "4. Start the frontend: cd frontend && npm run dev"
echo ""
echo "📖 For detailed instructions, see SETUP.md"
