## MODIFIED Requirements

### Requirement: Atomic Command Set
The system SHALL expose and document a single, stable set of atomic GitLab/git operations. Each atomic command SHALL perform exactly one identifiable action (e.g. list projects, get one pipeline, compare two branches, create MR, set auth config). The set SHALL be documented as the authoritative list (in specs and user-facing docs) and SHALL include: config auth; project list/get; pipeline list/get/latest/check-schedule; branch list/diff; mr list/create/merge; tag list/create/delete. Filtering and paging SHALL be expressed via existing flags (e.g. --status, --limit), not by splitting into additional subcommands.

#### Scenario: Atomic set is documented and discoverable
- **WHEN** a user or Agent consults the documented atomic command set (e.g. README or `gitlab-tools --help` and subcommand help)
- **THEN** the list of atomic commands SHALL match the implemented CLI subcommands for project, pipeline, branch, mr, tag, and config
- **AND** each entry SHALL have a short description and main parameters so that Agent can choose the right command

#### Scenario: Each atomic command has single responsibility
- **WHEN** an atomic command is invoked with valid arguments
- **THEN** it SHALL perform exactly one GitLab/git operation (e.g. list, get, create, merge, delete, or config auth) and SHALL not combine multiple unrelated operations in one call
- **AND** filtering and limits SHALL be done via flags (e.g. --status, --limit, --search) rather than extra subcommands

## ADDED Requirements

### Requirement: Global Config Path Flag
The system SHALL support a global flag (e.g. `-c`, `--config`) to specify the configuration file path. When specified, the tool SHALL use that file for reading and writing configuration for the duration of the command. The flag SHALL be documented in root command help so that users and Agent can discover it.

#### Scenario: Config path flag is documented
- **WHEN** a user or Agent runs `gitlab-tools --help` or consults CLI documentation
- **THEN** the global config path flag (e.g. `-c`, `--config`) SHALL be listed with a short description (e.g. config file path)
- **AND** so that users can specify a config file for that invocation (e.g. `gitlab-tools -c /path/to/config.yaml project list`)
