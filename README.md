# kkp-hackathon-car-park
Just for fun projects

## Backend (be/)
- Go 1.25 + Fiber v2 + PostgreSQL (GORM)
- Run full stack with Docker: `docker compose up --build`
- Run API: `cd be && go run ./cmd/api`

## Frontend (fe/)
- React + TypeScript + Vite
- Run UI: `cd fe && npm install && npm run dev`
- Frontend expects API at `http://localhost:8080`
- Optional override: `VITE_API_BASE_URL`

## Docker Compose Services
- `postgres`: PostgreSQL database on `localhost:5432`
- `be`: Backend API on `http://localhost:8080`
- `fe`: Frontend app on `http://localhost:5173`

## Implement-1 Delivered

### Slot status
- `free`
- `parked`
- `close`

### Slot type
- `VIP`
- `normal`

### Service
- Realtime parking space status via polling dashboard endpoint
- Parking fee calculation endpoint

### Kiosk Mode
- Show parking slot status
- Show summary and pricing information

### Management Mode
- Admin login system
- Admin can edit parking slot status

## Default Admin Credentials
- Username: `admin`
- Password: `admin123`
- Override with env vars:
	- `ADMIN_USERNAME`
	- `ADMIN_PASSWORD`

## Main API Endpoints
- `GET /health`
- `GET /kiosk/dashboard`
- `GET /kiosk/info`
- `GET /parking/fee?slotId=<id>&exitAt=<RFC3339 optional>`
- `POST /management/login`
- `GET /management/slots` (Bearer token)
- `PATCH /management/slots/:id/status` (Bearer token)
