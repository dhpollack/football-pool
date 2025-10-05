@default:
    @just -u -l --list-submodules

docker-compose-up:
    docker compose up --build -d

docker-compose-down:
    docker compose down

build-push-local:
    # This builds the amd image local and the arm image on the orangepi5-max-ubuntu
    just frontend build-push-local
    just backend build-push-local

_claude *args:
    claude {{ args }}

claude *args:
    just -E .env-claude _claude {{ args }}

helm-lint:
    helm lint helm/football-pool

mod backend 'backend/justfile'
mod frontend 'frontend/justfile'
