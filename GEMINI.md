# Football Pool Webapp and Backend

## Summary

We are creating a webapp for a football pool.  The technical design details exist as markdown files in the `design` directory.

## Development Plan

You have created a development plan in the file `design/development_plan.md`.  While implementing this project, you should keep track of our progress there.

## Commands

We can use the `just` recipes in the terminal, but in the dockerfile we should build and run the commands directly so we don't need `just` in the docker image

### Backend

- `just backend format`: format the code, this should be used before linting
- `just backend lint`: lint the backend app
- `just backend build`: build the golang server api app
- `just backend run`: run the backend api server

### Frontend

- `just frontend format`: format the code with biome
- `just frontend lint`: lint code with biome
