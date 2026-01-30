## ADDED Requirements

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output (e.g. `--json`) for all commands in this capability (project list, project get). When the designated flag is set, the command SHALL output the primary result as JSON to stdout with a stable structure, as specified by the agent-interface capability.

#### Scenario: JSON output for project commands
- **WHEN** user or Agent runs `gitlab-tools project list` or `gitlab-tools project get <id>` with the machine-output flag (e.g. `--json`)
- **THEN** the system SHALL print valid JSON to stdout with the key fields of the result (e.g. id, path, visibility, web_url for projects)
- **AND** the output SHALL be parseable by scripts and Agent for chaining or decision-making

### Requirement: List Projects with Filters
The system SHALL provide `project list` as an atomic command that returns projects accessible to the user. Filtering SHALL be done via flags only (e.g. `--owned`, `--search`, `--match`, `--limit`, `--has-schedule`); no separate subcommands for filter types. Output SHALL include project id and path so that pipeline/branch/mr/tag commands can use either as input.

#### Scenario: List with optional filters
- **WHEN** user or Agent runs `gitlab-tools project list` with optional flags (e.g. `--search "backend"`, `--limit 50`, `--owned`)
- **THEN** the system SHALL return only projects matching the filters
- **AND** each project in the output SHALL include id and path (or path_with_namespace) for use as project identifier in other commands

#### Scenario: Empty list is valid result
- **WHEN** no projects match the filters
- **THEN** the system SHALL return an empty list and exit code 0
- **AND** SHALL NOT treat empty as an error

### Requirement: Project List Filter Semantics
The semantics of `--search` and `--match` SHALL be documented: `--search` SHALL match project name or description by substring (typically case-insensitive); `--match` SHALL match project path or name by regular expression. So that Agent can choose the right filter for "find by name" vs "find by path pattern".

#### Scenario: Search and match documented
- **WHEN** user or Agent consults docs or help for `project list`
- **THEN** the difference between `--search` and `--match` SHALL be stated (substring vs regex, and which fields are matched)
- **AND** so that Agent can construct correct filters (e.g. `--match "^group/.*"` for path prefix)

### Requirement: Project List Quiet Mode
The system SHALL support `project list` with optional `--quiet` that outputs only project identifiers (e.g. project path or name, one per line) with no other fields. The output shape SHALL be documented so that Agent can pipe or parse a minimal list for chaining (e.g. feed to pipeline list).

#### Scenario: Quiet output shape
- **WHEN** user or Agent runs `gitlab-tools project list [filters] --quiet`
- **THEN** the system SHALL output only one project identifier per line (path or name as documented)
- **AND** no headers or extra fields SHALL be printed so that output is suitable for scripting

### Requirement: Project Get Output for Composition
The system SHALL ensure `project get <id|path>` output includes at least: id (numeric), path (or path_with_namespace), default_branch, web_url. So that Agent or scripts can chain to branch list, pipeline latest, or mr list using the same identifier.

#### Scenario: Get output includes chaining fields
- **WHEN** user or Agent runs `gitlab-tools project get <id|path>` with machine-readable output
- **THEN** the output SHALL include id, path, default_branch, and web_url
- **AND** the same id or path MAY be passed to `branch list`, `pipeline latest`, `mr list`, `tag list` as the project argument

## REMOVED Requirements

### Requirement: Architecture Pattern
**Reason**: Specs focus on user- and Agent-visible behavior; code layout and package structure are implementation concerns and belong in project conventions or design docs, not in a domain capability spec.
**Migration**: Document code organization in `openspec/project.md` or a dedicated development-conventions section; remove this requirement from project-management spec when archiving this change.
