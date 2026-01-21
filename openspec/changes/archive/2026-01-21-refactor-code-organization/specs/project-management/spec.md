## ADDED Requirements
### Requirement: Architecture Pattern
The system SHALL follow a modular package structure for code organization.

#### Scenario: Code organization structure
- **WHEN** the codebase is organized
- **THEN** the system SHALL use a standard Go project layout with `cmd/` and `internal/` directories
- **AND** the system SHALL separate functionality into distinct packages (pipeline, project, branch, mr, tag)
- **AND** the system SHALL extract shared logic (client creation, configuration management) into common packages
- **AND** each functional module SHALL be organized in its own package under `internal/`
- **AND** the main entry point SHALL be located in `cmd/gitlab-tools/main.go`

#### Scenario: Package structure
- **WHEN** code is organized into packages
- **THEN** the system SHALL place GitLab client creation logic in `internal/client`
- **AND** the system SHALL place configuration management logic in `internal/config`
- **AND** the system SHALL place each command group (pipeline, project, branch, mr, tag) in its own package under `internal/`
- **AND** each command group package SHALL contain command definitions, command handlers, and output formatting functions

