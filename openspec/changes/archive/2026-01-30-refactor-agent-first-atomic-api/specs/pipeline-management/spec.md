## ADDED Requirements

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output (e.g. `--json`) for all commands in this capability (pipeline list, pipeline get, pipeline latest, pipeline check-schedule). When the designated flag is set, the command SHALL output the primary result as JSON to stdout with a stable structure, as specified by the agent-interface capability.

#### Scenario: JSON output for pipeline commands
- **WHEN** user or Agent runs a pipeline command (e.g. `gitlab-tools pipeline list <project>`, `gitlab-tools pipeline latest <project> <ref>`) with the machine-output flag (e.g. `--json`)
- **THEN** the system SHALL print valid JSON to stdout with the key fields of the result (e.g. id, status, ref, web_url, created_at for pipelines)
- **AND** the output SHALL be parseable by scripts and Agent for chaining or decision-making

### Requirement: Pipeline List Filters and Output
The system SHALL support `pipeline list <project> [project ...]` with optional flags `--status`, `--ref`, `--limit`. Status SHALL accept only documented values (e.g. running, pending, success, failed, canceled, skipped, created, manual). Output SHALL include per-item: pipeline id, status, ref, sha (or commit short_sha), web_url, created_at; and project identification when multiple projects are listed.

#### Scenario: List with status and limit
- **WHEN** user or Agent runs `gitlab-tools pipeline list <project> --status success --limit 10`
- **THEN** the system SHALL return only pipelines with the given status, limited to the given count
- **AND** each item SHALL include id, status, ref, web_url so that `pipeline get <project> <id>` can be chained

#### Scenario: Invalid status returns usage error
- **WHEN** user or Agent passes an invalid value to `--status`
- **THEN** the system SHALL emit a usage or validation error and exit with code 2
- **AND** SHALL NOT call the GitLab API with the invalid value

#### Scenario: Multiple projects each have project identification
- **WHEN** user or Agent runs `gitlab-tools pipeline list <project1> <project2> ...` with multiple projects
- **THEN** the output SHALL be grouped by project and each group SHALL include project id or path so that pipeline ids can be associated with the correct project for chaining (e.g. pipeline get <project> <id>)

### Requirement: Pipeline List Ordering
Pipelines in `pipeline list` output SHALL be ordered in a documented way (e.g. by created_at descending, newest first). This order SHALL be documented so that Agent can rely on the first item as "most recent" when limit is used.

#### Scenario: Order documented
- **WHEN** user or Agent consults docs or help for `pipeline list`
- **THEN** the sort order (e.g. newest first) SHALL be stated
- **AND** so that "first pipeline" or "first N pipelines" has a well-defined meaning

### Requirement: Pipeline Get and Latest Output for Composition
The system SHALL ensure `pipeline get <project> <pipeline-id>` and `pipeline latest <project> <ref>` output includes: id, status, ref, sha (or short_sha), web_url, created_at. So that Agent can decide next action (e.g. open URL, retry, or chain to mr list by ref).

#### Scenario: Get and latest output include chaining fields
- **WHEN** user or Agent runs `pipeline get <project> <id>` or `pipeline latest <project> <ref>` with machine-readable output
- **THEN** the output SHALL include id, status, ref, sha, web_url, created_at
- **AND** ref SHALL be usable as branch/tag name for `branch diff` or `mr list --target-branch` in composition

### Requirement: Get Single Pipeline
The system SHALL provide `pipeline get <project> <pipeline-id>` as an atomic command that returns one pipeline by project and pipeline ID. Output SHALL include id, status, ref, sha, web_url, created_at. When the pipeline or project does not exist or is not accessible, the system SHALL return an error and exit with code 1.

#### Scenario: Get pipeline by project and id
- **WHEN** user or Agent runs `gitlab-tools pipeline get <project> <pipeline-id>` with valid project (id or path) and pipeline id
- **THEN** the system SHALL fetch and return that pipeline
- **AND** the output SHALL include id, status, ref, sha, web_url, created_at for composition or display

#### Scenario: Pipeline or project not found
- **WHEN** the pipeline id or project does not exist or the user has no access
- **THEN** the system SHALL return an error message and exit with code 1
- **AND** SHALL NOT exit with code 0

### Requirement: Latest Pipeline by Ref
The system SHALL provide `pipeline latest <project> <ref>` as an atomic command that returns the most recent pipeline for the given ref (branch or tag name). Output SHALL include id, status, ref, sha, web_url, created_at. When no pipeline exists for the ref, the system SHALL return an appropriate message and exit with code 1 so Agent can distinguish "no pipeline" from success.

#### Scenario: Latest by branch or tag
- **WHEN** user or Agent runs `gitlab-tools pipeline latest <project> <ref>` with valid project and ref (e.g. main, develop, v1.0)
- **THEN** the system SHALL return the latest pipeline for that ref
- **AND** the output SHALL include id, status, ref, sha, web_url, created_at

#### Scenario: No pipeline for ref
- **WHEN** no pipeline exists for the given ref
- **THEN** the system SHALL return an appropriate message and exit with code 1
- **AND** so that Agent can branch (e.g. trigger pipeline or report)

### Requirement: Check Schedule
The system SHALL provide `pipeline check-schedule <project> [schedule-id]` as an atomic command that checks whether the most recent scheduled pipeline (or the given schedule) succeeded. Output SHALL indicate success or failure; when machine-readable output is used, the result SHALL be parseable so that Agent can decide next action (e.g. alert, retry).

#### Scenario: Check schedule result is clear
- **WHEN** user or Agent runs `gitlab-tools pipeline check-schedule <project>` or with optional schedule-id
- **THEN** the system SHALL report whether the relevant scheduled pipeline succeeded or failed (or not found)
- **AND** exit code SHALL be 0 on success and 1 on failure or missing schedule, per agent-interface exit code contract

#### Scenario: No schedule configured or invalid schedule-id
- **WHEN** the project has no pipeline schedule or the given schedule-id does not exist
- **THEN** the system SHALL return an appropriate message and exit with code 1
- **AND** the message SHALL indicate the reason (e.g. no schedule, schedule not found) so that Agent can decide next action

### Requirement: Allowed Pipeline Status Values
The allowed values for `pipeline list --status` SHALL be documented (e.g. in `--help` or docs) and SHALL be validated before calling the API. Recommended set: running, pending, success, failed, canceled, skipped, created, manual. Invalid values SHALL result in exit code 2 per agent-interface.

#### Scenario: Allowed status values are documented
- **WHEN** user or Agent runs `gitlab-tools pipeline list --help` or consults docs
- **THEN** the allowed values for `--status` SHALL be listed or linked
- **AND** so that Agent can choose a valid value without trial and error
