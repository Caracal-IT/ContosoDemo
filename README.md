# Contoso REST API

A starter REST API project using Gin.

## Getting Started

1. Install backend dependencies:
   ```
   go mod tidy
   ```

2. **Build the frontend before running the backend:**
   ```
   cd frontend
   npm install
   npm run build
   cd ..
   ```
   The built files will be output to `public`.

3. Run the server:
   ```
   go run main.go
   ```

4. Test the API:
   ```
   curl http://localhost:8080/api/ping
   ```

## Automated Build

To ensure the frontend is always built before running the backend, you can use a script:

**On Unix/macOS:**
```sh
#!/bin/sh
cd frontend
npm install
npm run build
cd ..
go run main.go
```

**On Windows (PowerShell):**
```powershell
cd frontend
npm install
npm run build
cd ..
go run main.go
```

Or add a Makefile or npm script to automate this process.

## Structure

- `main.go` - Entry point
- `routes/` - Route definitions
- `controllers/` - Request handlers
- `models/` - Data models
- `frontend/` - Vue.js web frontend (Quasar, Vite)
- `public/` - Built frontend files (served by Go)

## ELK Stack
git clone https://github.com/deviantony/docker-elk.git
