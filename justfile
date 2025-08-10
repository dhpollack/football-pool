@default:
  @just -u -l --list-submodules

docker-compose-up:
  docker compose up -d

docker-compose-down:
  docker compose down

mod backend 'backend/justfile'

mod frontend 'frontend/justfile'

