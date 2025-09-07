@echo off
echo 🚀 Setting up AI Study Partner...

REM Check if required tools are installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install it first.
    pause
    exit /b 1
) else (
    echo ✅ Go is installed
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ Node.js is not installed. Please install it first.
    pause
    exit /b 1
) else (
    echo ✅ Node.js is installed
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo ❌ npm is not installed. Please install it first.
    pause
    exit /b 1
) else (
    echo ✅ npm is installed
)

REM Setup backend
echo 🔧 Setting up backend...
cd backend

REM Copy environment file if it doesn't exist
if not exist .env (
    copy env.example .env
    echo 📝 Created .env file. Please edit it with your configuration.
)

REM Install Go dependencies
echo 📦 Installing Go dependencies...
go mod tidy

cd ..

REM Setup frontend
echo 🎨 Setting up frontend...
cd frontend

REM Copy environment file if it doesn't exist
if not exist .env.local (
    copy env.example .env.local
    echo 📝 Created .env.local file. Please edit it with your configuration.
)

REM Install Node.js dependencies
echo 📦 Installing Node.js dependencies...
npm install

cd ..

echo ✅ Setup complete!
echo.
echo 📋 Next steps:
echo 1. Set up your database (see SETUP.md for details)
echo 2. Configure your .env files
echo 3. Start the backend: cd backend ^&^& go run cmd/main.go
echo 4. Start the frontend: cd frontend ^&^& npm run dev
echo.
echo 📖 For detailed instructions, see SETUP.md
pause
