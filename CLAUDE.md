# Football Pool Webapp and Backend

## Summary

We are creating a webapp for a football pool.  The technical design details exist as markdown files in the `design` directory.

## Quality Controls

- All code should be properly formatted and linted
- All unit tests should pass
- All integration tests should pass
- Never edit generated code files

## Available Tools

- Github MCP server: we are going to use github issues to track major issues and progress
  - repo: github.com/dhpollack/football-pool
- just: this is a command runner that can simplify running commands
- code generation tools
  - `oapi-codegen`: (backend) generate data models from openapi spec
  - `orval`: (frontend) generate react query client and data models from openapi spec

## Commands

### Backend

- `just backend format`: format the code, this should be used before linting
- `just backend lint`: lint the backend app
- `just backend generate`: generate data models from openapi spec
- `just backend build`: build the golang server api app
- `just backend run`: run the backend api server

### Frontend

For the frontend, we also have a just file, but you should use the npm scripts.  Importantly, *NEVER* use `npx` to run tools.

- `npm run format`: format files with biome
- `npm run lint`: lint files with biome
- `npm run generate`: generate the frontend react-query client and data models and msw mocks
- `npm run build`: build production code and check for typescript errors
- `npm run test`: run the unit tests
- `npm run test:e2e`: run e2e tests on all platforms (only should be used once all tests are passing)
- `npm run test:e2e:chrome`: run the e2e tests with the headless chrome runner.  This is much faster and should be used for most tests
