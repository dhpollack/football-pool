@default:
    @just -u -l --list-submodules

docker-compose-up:
    docker compose up --build -d

docker-compose-down:
    docker compose down

_claude:
    claude

claude:
    just -E .env-claude _claude

mod backend 'backend/justfile'
mod frontend 'frontend/justfile'
