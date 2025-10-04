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

## Helm Chart

This project includes a Helm chart for deploying the application to a Kubernetes cluster.

### Prerequisites

- Helm v3+
- A Kubernetes cluster
- `kubectl` configured to connect to your cluster

### Installation

Install the Helm chart from the public OCI registry:

```bash
helm install my-release oci://ghcr.io/dhpollack/football-pool --version 0.2.1
```

### Configuration

You can customize the installation by creating a `values.yaml` file and passing it with the `--values` flag.

For example, to enable the managed PostgreSQL cluster:

**`my-values.yaml`:**
```yaml
postgresql:
  enabled: true
  storage:
    size: "5Gi"
```

**Install command:**
```bash
helm install my-release oci://ghcr.io/dhpollack/football-pool --version 0.2.1 -f my-values.yaml
```

See the `helm/football-pool/values.yaml` file for all available configuration options.

### Terraform Usage

You can also deploy the Helm chart using Terraform:

```hcl
resource "helm_release" "football_pool" {
  name       = "football-pool"
  repository = "oci://ghcr.io/dhpollack"
  chart      = "football-pool"
  version    = "0.2.1"

  values = [
    file("${path.module}/values.yaml")
  ]
}
```

## Notes

We can get data about the games from espn api.  See: https://gist.github.com/nntrn/ee26cb2a0716de0947a0a4e9a157bc1c
