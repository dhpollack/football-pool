# Football Pool Webapp and Backend

This repository contains the web application and backend for a football pool game.

## Project Structure

- `backend/`: Contains the Go backend application.
- `frontend/`: Contains the React frontend application.
- `design/`: Contains the design documents for the project.

## Backend Setup and Running

### Prerequisites

- Go (version 1.16 or higher)

### Setup

1. Navigate to the `backend` directory:
   ```bash
   cd backend
   ```
2. Initialize Go modules and download dependencies:
   ```bash
   go mod tidy
   ```

### Running

To run the backend server:

```bash
cd backend
go run cmd/main.go
```

The backend server will start on `http://localhost:8080`.

## Docker Compose

To bring up the entire application using Docker Compose:

```bash
just docker-compose-up
```

To bring down the application:

```bash
just docker-compose-down
```

## Just Commands

This project uses `just` for task automation. You can see available commands by running:

```bash
just
```

To run commands within specific modules (backend or frontend):

```bash
just mod backend <command>
just mod frontend <command>
```

For example, to lint the backend code:

```bash
just mod backend lint
```

## Frontend Setup and Running

### Prerequisites

- Node.js (LTS version recommended)
- npm (comes with Node.js)

### Setup

1. Navigate to the `frontend` directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```

### Running

To run the frontend development server:

```bash
cd frontend
npm run dev
```

The frontend application will open in your browser, usually at `http://localhost:5173`.

