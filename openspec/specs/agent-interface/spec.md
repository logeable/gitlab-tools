# agent-interface Specification

## Purpose
TBD - created by archiving change refactor-agent-first-atomic-api. Update Purpose after archive.
## Requirements
### Requirement: Atomic Command Set
The system SHALL expose and document a single, stable set of atomic GitLab/git operations. Each atomic command SHALL perform exactly one identifiable action (e.g. list projects, get one pipeline, compare two branches, create MR). The set SHALL be documented as the authoritative list (in specs and user-facing docs) and SHALL include: project list/get; pipeline list/get/latest/check-schedule; branch list/diff; mr list/create/merge; tag list/create/delete. Filtering and paging SHALL be expressed via existing flags (e.g. --status, --limit), not by splitting into additional subcommands.

#### Scenario: Atomic set is documented and discoverable
- **WHEN** a user or Agent consults the documented atomic command set (e.g. README or `gitlab-tools --help` and subcommand help)
- **THEN** the list of atomic commands SHALL match the implemented CLI subcommands for project, pipeline, branch, mr, and tag
- **AND** each entry SHALL have a short description and main parameters so that Agent can choose the right command

#### Scenario: Each atomic command has single responsibility
- **WHEN** an atomic command is invoked with valid arguments
- **THEN** it SHALL perform exactly one GitLab/git operation (e.g. list, get, create, merge, delete) and SHALL not combine multiple unrelated operations in one call
- **AND** filtering and limits SHALL be done via flags (e.g. --status, --limit, --search) rather than extra subcommands

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output for all atomic commands. When the designated flag (e.g. `--json`) is set, the command SHALL output the primary result as JSON to stdout and SHALL use a stable, documented structure (object or array) with fields sufficient for Agent or script parsing (e.g. id, status, url, timestamps). Errors SHALL be reported in a consistent way (e.g. stderr and/or structured JSON error field) so that callers can distinguish success from failure.

#### Scenario: JSON output when flag is set
- **WHEN** user or Agent runs an atomic command with the machine-output flag (e.g. `gitlab-tools project list --json`)
- **THEN** the system SHALL print valid JSON to stdout representing the command result
- **AND** the JSON SHALL include the key fields that are shown in human-readable output (e.g. project id, path, visibility; pipeline id, status, ref, web_url)

#### Scenario: Human-readable default
- **WHEN** the machine-output flag is not set
- **THEN** the system SHALL behave as today: human-readable output and messages in the configured language (e.g. Chinese), with no breaking change to existing scripts or users

### Requirement: Exit Code Contract
The system SHALL use a stable exit code contract for all atomic commands: 0 for success (including empty results), 1 for business or API errors (e.g. project not found, permission denied, MR already merged), 2 for usage errors (e.g. missing required argument, invalid flag). The root command and all subcommands SHALL adhere to this contract so that Agent and shell scripts can branch on exit code.

#### Scenario: Success exit code
- **WHEN** an atomic command completes successfully (including returning an empty list)
- **THEN** the process SHALL exit with code 0

#### Scenario: Business or API error exit code
- **WHEN** an atomic command fails due to GitLab API or business rules (e.g. project not found, insufficient permissions, merge conflict)
- **THEN** the process SHALL exit with code 1
- **AND** an appropriate error message SHALL be emitted to stderr (or structured error when --json is used)

#### Scenario: Usage error exit code
- **WHEN** an atomic command is invoked with missing required arguments or invalid flags
- **THEN** the process SHALL exit with code 2
- **AND** usage or validation error information SHALL be emitted to stderr

### Requirement: Composition and Workflow Discoverability
The system SHALL document composition patterns and workflow recipes so that Agent and users can combine atomic commands to achieve common scenarios (e.g. discover project → list branches → diff → create MR; or latest pipeline by branch and status). Documentation SHALL be available in the repository (e.g. README or docs) and SHALL be reflected in the Agent Skill (e.g. SKILL.md) with a clear split between "atomic commands" and "recommended combinations". The tool MAY provide a discoverability command (e.g. `gitlab-tools capabilities`) that outputs the list of atomic capabilities with short descriptions.

#### Scenario: Workflow documentation exists
- **WHEN** a user or Agent looks for how to combine commands (e.g. "find project then create MR")
- **THEN** the repository SHALL provide at least one documented section (e.g. "Workflow and composition examples") with at least three concrete scenarios
- **AND** each scenario SHALL describe the sequence of atomic commands and, where useful, note that output can be parsed with --json for chaining

#### Scenario: Skill distinguishes atoms and combinations
- **WHEN** an Agent uses the provided Skill (e.g. SKILL.md)
- **THEN** the Skill SHALL list atomic commands in a table or structured form (name, main args, typical use)
- **AND** the Skill SHALL include at least one or two recommended combination examples (e.g. one-liner or step list) and SHALL mention --json for programmatic use

### Requirement: Unified Project Identifier
All commands that operate on a single GitLab project SHALL accept the project in one of two forms: numeric project ID or project path (PathWithNamespace, e.g. `group/subgroup/project`). The system SHALL resolve either form to the same project. Output that refers to a project SHALL include both `id` (numeric) and `path` (or equivalent) so that downstream commands may use either form.

#### Scenario: Accept numeric ID or path
- **WHEN** user or Agent invokes any project-scoped command (e.g. pipeline list, branch list, mr list, tag list) with either a numeric ID or a project path
- **THEN** the system SHALL treat both as the same project and SHALL produce the same logical result
- **AND** resolution SHALL be consistent with GitLab API (get project by ID or by path)

#### Scenario: Output includes both id and path for chaining
- **WHEN** a command returns project-related data (e.g. project list, project get, or grouped results by project)
- **THEN** each project reference in the output SHALL include at least numeric id and path (or path_with_namespace) when machine-readable output is used
- **AND** so that the next command in a composition MAY use either id or path as the project argument

### Requirement: Composable Output Fields
Output of list/get commands SHALL include fields that are sufficient for composition: identifiers (id, path, iid, ref, sha as applicable), web_url for human or Agent navigation, and timestamps where relevant. Field names SHALL be stable and documented (e.g. in --help or docs) so that Agent and scripts can parse and chain without guessing.

#### Scenario: List output supports chaining
- **WHEN** user or Agent runs a list command (e.g. project list, pipeline list, mr list) with machine-readable output
- **THEN** each item SHALL include the primary identifier(s) needed as input to another atomic command (e.g. project id/path for pipeline list; pipeline id for pipeline get; mr iid for mr merge)
- **AND** web_url SHALL be included when applicable for traceability

#### Scenario: Single-item output supports chaining
- **WHEN** user or Agent runs a get/create/merge/delete command with machine-readable output
- **THEN** the result SHALL include the created or affected resource identifier and web_url when applicable
- **AND** so that follow-up commands or logging can use the same identifiers

### Requirement: Empty Result Is Success
List and get commands SHALL treat "no items found" or "empty list" as success, not as an error. The process SHALL exit with code 0 and SHALL output an empty array (or equivalent) when machine-readable output is used, so that Agent and scripts can distinguish "no data" from "failure" without relying on exit code.

#### Scenario: Empty list returns exit 0
- **WHEN** a list command returns zero items (e.g. no projects match filters, no pipelines, no MRs)
- **THEN** the process SHALL exit with code 0
- **AND** when machine-readable output is used, the result SHALL be an empty array or documented "empty" structure, not an error payload

### Requirement: Pagination and Limit
List commands SHALL support a limit on the number of items returned (e.g. `--limit`). The default and maximum limit SHALL be documented (e.g. in `--help` or docs) so that Agent can bound response size and avoid unbounded output. When limit is not specified, a reasonable default SHALL be used.

#### Scenario: Limit is documented and applied
- **WHEN** user or Agent runs a list command with or without `--limit <n>`
- **THEN** the system SHALL return at most n items when limit is specified
- **AND** the default limit and maximum (if any) SHALL be documented so Agent can choose appropriately

### Requirement: JSON Field Naming
When machine-readable output (e.g. JSON) is used, field names SHALL use snake_case (e.g. `web_url`, `source_branch`, `created_at`) for consistency with GitLab API and common tooling. This SHALL be documented so that Agent and scripts can rely on stable field names.

#### Scenario: JSON uses snake_case fields
- **WHEN** user or Agent runs any command with machine-readable output
- **THEN** the JSON SHALL use snake_case for all field names (e.g. web_url, not webUrl)
- **AND** the same field names SHALL be used across list/get/create/merge/delete for the same concept

### Requirement: Documented Enum Values
Where a command accepts an enumerated value (e.g. pipeline status, MR state), the allowed values SHALL be documented (e.g. in `--help` or docs) and SHALL be validated before calling the API. Invalid values SHALL result in exit code 2 (usage error). This allows Agent to discover valid options without trial and error.

#### Scenario: Invalid enum returns usage error
- **WHEN** user or Agent passes an invalid value for a documented enum (e.g. `--status invalid`, `--state invalid`)
- **THEN** the system SHALL emit a usage or validation error and exit with code 2
- **AND** the error message SHOULD indicate the allowed values or where to find them

### Requirement: Required and Optional Arguments Documented
For each atomic command, required positional arguments and optional flags SHALL be documented (e.g. in `--help` or docs). Required arguments SHALL be clearly distinguished from optional ones so that Agent and scripts know the minimal invocation and can discover options without guessing.

#### Scenario: Help shows required vs optional
- **WHEN** user or Agent runs `gitlab-tools <subcommand> --help`
- **THEN** the usage line SHALL show required positional arguments and optional flags
- **AND** so that Agent can construct a valid command with only required args and add flags as needed

### Requirement: Structured Error Output When Machine-Readable
When the machine-output flag (e.g. `--json`) is set and an error occurs (exit 1 or 2), the system SHALL emit a structured error to stderr (e.g. a JSON object with an error message) so that Agent can parse errors consistently. The exact field names SHALL be documented. Exit code SHALL still be 1 or 2 per Exit Code Contract so that scripts can branch on code without parsing stderr.

#### Scenario: JSON error on stderr when --json and failure
- **WHEN** user or Agent runs a command with `--json` and the command fails (business error or usage error)
- **THEN** the system SHALL output a parseable error payload to stderr (e.g. `{"error": "message"}` or documented format)
- **AND** exit code SHALL still be 1 or 2 per Exit Code Contract so that scripts can branch on code without parsing stderr

### Requirement: Terminology for Ref and Branches
The following terms SHALL be used consistently so that Agent and docs can align: "ref" denotes a branch name or tag name usable in commands such as `pipeline latest <project> <ref>` and `tag list`; "source_branch" and "target_branch" denote branch names in MR context (e.g. `mr create <project> <source_branch> <target_branch>`). Documentation and machine-readable field names SHALL use these terms where applicable.

#### Scenario: Consistent ref and branch terms
- **WHEN** user or Agent reads docs or parses JSON output
- **THEN** "ref" SHALL refer to a branch or tag name in pipeline/tag context
- **AND** "source_branch" / "target_branch" SHALL refer to MR source and target branches; the same branch name MAY appear as "ref" in pipeline context and as "source_branch" or "target_branch" in MR context

### Requirement: List Output Ordering Documented
When a list command returns items in a well-defined order (e.g. newest first by created_at), that order SHALL be documented (e.g. in `--help` or docs) so that Agent can rely on "first" or "last" item (e.g. "first pipeline" = most recent) without guessing.

#### Scenario: Order is documented
- **WHEN** user or Agent consults docs or help for a list command (e.g. pipeline list, mr list, tag list)
- **THEN** the documentation SHALL state the sort order when it is deterministic (e.g. pipelines by created_at descending)
- **AND** so that Agent can use the first or last element for "latest" or "oldest" without additional filtering

### Requirement: Mutation Clear Errors
Mutations (create, merge, delete) SHALL return a clear error and exit with code 1 when the operation is not applicable: e.g. resource already exists (duplicate tag, MR already exists), resource already in target state (MR already merged), or precondition failed (conflicts, invalid ref). The error message SHALL indicate the reason so that Agent can decide next action (e.g. skip, retry, or report).

#### Scenario: Duplicate or already-done returns clear error
- **WHEN** user or Agent attempts a mutation that is not applicable (e.g. tag create with existing name, mr merge on already-merged MR)
- **THEN** the system SHALL return an error message and exit with code 1
- **AND** the message SHALL indicate the reason (e.g. "tag already exists", "merge request already merged") so that Agent can branch without parsing API-specific codes

### Requirement: Output Schema Stability
JSON field names and the top-level structure (object vs array) for machine-readable output SHALL remain stable across minor versions of the tool. Breaking changes to field names or structure SHALL be documented (e.g. in changelog or release notes) and SHOULD be avoided or versioned (e.g. via a format version field) so that Agent and scripts can rely on stable parsing.

#### Scenario: Same command same structure
- **WHEN** user or Agent runs the same command with the same arguments across a minor version upgrade
- **THEN** the JSON output SHALL retain the same field names and logical structure unless a breaking change is explicitly documented
- **AND** so that existing scripts and Agent prompts do not break silently

