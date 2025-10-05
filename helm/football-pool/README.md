# football-pool helm chart

A helm chart the football-pool

## Prerequisites

- Helm v3+
- A Kubernetes cluster
- `kubectl` configured to connect to your cluster

## Installation

Install the Helm chart from the public OCI registry:

```bash
helm install my-release oci://ghcr.io/dhpollack/football-pool --version 0.2.1
```

## Configuration

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

## Terraform Usage

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

## Contributing

Make sure the helm chart lints correctly with the following command:

```shell
helm lint .
```

