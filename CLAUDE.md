# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build Commands
- `task build` - Build the application binary to `bin/jelly`
- `task gen` - Generate API code from OpenAPI spec and mocks for testing

### Code Generation
- `oapi-codegen --config=config/oapi-codegen.yaml config/api.yaml` - Generate API server code
- `mockery --config=config/mockery.yaml` - Generate test mocks

### Testing
- `task test` - Run all tests
- `task test-verbose` - Run tests with verbose output
- `task test-cov` - Run tests with coverage
- `task test-cov-out` - Run tests with coverage output to HTML

### Development
- `task run-local` - Run the application locally with `ENVIRONMENT=local`

## Architecture

### API Endpoints
- `GET /health` - Health check endpoint
- `POST /photo` - Photo upload with optional caption and tags (multipart/form-data)

### Development Notes
- The `pkg/api.go` file is currently empty and needs implementation
- The project uses Task runner instead of Make
- OpenAPI spec generates server code to `pkg/api/gen/api.gen.go`
- Environment-based configuration (local/dev enables pprof on port 6060)

### Generated Code
- API server code is generated from the OpenAPI spec
- Test mocks are generated using mockery
- Always run `task gen` before building to ensure generated code is up to date
- OpenAPI 3 specification generated code is stored in @pkg/api/gen/ 

### Development Constraints
- Do not edit files in folders named "gen"

### Naming and Code Generation
- Handler functions for endpoints must match the codegen and OpenAPI specification names with public naming schema

### Documentation
- Documentation for dependencies is stored in the `claude` folder if it exists
- ALWAYS read the README.md file for a list of paths to the relevant docs

