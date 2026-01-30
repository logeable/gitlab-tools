## ADDED Requirements

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output (e.g. `--json`) for all commands in this capability (mr list, mr create, mr merge). When the designated flag is set, the command SHALL output the primary result as JSON to stdout with a stable structure, as specified by the agent-interface capability.

#### Scenario: JSON output for MR commands
- **WHEN** user or Agent runs `gitlab-tools mr list <project>`, `gitlab-tools mr create ...`, or `gitlab-tools mr merge <project> <iid>` with the machine-output flag (e.g. `--json`)
- **THEN** the system SHALL print valid JSON to stdout with the key fields of the result (e.g. iid, title, source_branch, target_branch, web_url for MRs)
- **AND** the output SHALL be parseable by scripts and Agent for chaining or decision-making

### Requirement: MR List Filters and Output
The system SHALL support `mr list <project>` with optional flags `--target-branch`, `--state` (e.g. opened, closed, merged). Output SHALL include per MR: iid, title, source_branch, target_branch, state, web_url, and optionally pipeline status. So that Agent can select an MR and chain to `mr merge <project> <iid>`.

#### Scenario: List output includes iid and branches for chaining
- **WHEN** user or Agent runs `gitlab-tools mr list <project>` with machine-readable output
- **THEN** each MR SHALL include iid, source_branch, target_branch, web_url
- **AND** iid SHALL be usable as the second argument to `mr merge <project> <iid>`

### Requirement: Merge Merge Request Output for Composition
The system SHALL provide `mr merge <project> <iid>` as an atomic command with optional `--delete-source-branch`, `--merge-commit-message`. Output SHALL include merge result and web_url. The system SHALL accept project as id or path per agent-interface. When the MR is already merged, has conflicts, or the user lacks permission, the system SHALL return an error and exit with code 1 so Agent can distinguish failure from success.

#### Scenario: Merge output supports traceability
- **WHEN** user or Agent runs `gitlab-tools mr merge <project> <iid>` with machine-readable output and merge succeeds
- **THEN** the output SHALL include merge result and web_url
- **AND** so that Agent or scripts can log or open the merged MR

#### Scenario: Merge failure returns exit 1
- **WHEN** the MR is already merged, has conflicts, or the user lacks permission to merge
- **THEN** the system SHALL return an appropriate error message and exit with code 1
- **AND** the error message SHALL indicate the reason (e.g. already merged, conflicts) so Agent can decide next action

### Requirement: Create Merge Request
The system SHALL provide `mr create <project> <source_branch> <target_branch>` as an atomic command with optional `--title`, `--description`. It SHALL accept project as id or path per agent-interface. Output SHALL include iid and web_url. When an MR from source to target already exists or creation fails (e.g. permission, invalid branch), the system SHALL return an error and exit with code 1.

#### Scenario: Create MR by project and branches
- **WHEN** user or Agent runs `gitlab-tools mr create <project> <source> <target>` with valid project and branch names
- **THEN** the system SHALL create one Merge Request from source_branch to target_branch
- **AND** the output SHALL include iid and web_url for chaining to mr merge or display

#### Scenario: Create with title and description
- **WHEN** user or Agent runs `gitlab-tools mr create <project> <source> <target> --title "..." --description "..."`
- **THEN** the system SHALL create the MR with the given title and description
- **AND** the output SHALL include iid and web_url

#### Scenario: MR already exists or creation fails
- **WHEN** an MR from source to target already exists, or branches are invalid, or user lacks permission
- **THEN** the system SHALL return an appropriate error message and exit with code 1
- **AND** the error message SHALL indicate the reason (e.g. MR exists, branch not found) so Agent can decide next action

### Requirement: MR List State Filter
The system SHALL support `mr list <project>` with optional `--state` (e.g. opened, closed, merged) so that Agent can list not only open MRs but also closed or merged ones for auditing or composition. Allowed state values SHALL be documented and validated; invalid values SHALL result in exit code 2.

#### Scenario: List by state
- **WHEN** user or Agent runs `gitlab-tools mr list <project> --state opened` or `--state merged`
- **THEN** the system SHALL return only MRs in the given state
- **AND** each item SHALL include iid, source_branch, target_branch, state, web_url for composition or reporting

### Requirement: MR List Default State
When `--state` is omitted for `mr list <project>`, the default state SHALL be documented (e.g. opened). So that Agent knows that `mr list <project>` without flags returns open MRs and can explicitly pass `--state closed` or `--state merged` when needed.

#### Scenario: Default state documented
- **WHEN** user or Agent consults docs or help for `mr list`
- **THEN** the default value for `--state` when omitted SHALL be stated (e.g. "default: opened")
- **AND** so that Agent can rely on "list MRs" without --state meaning "list open MRs"

### Requirement: Allowed MR State Values
The allowed values for `mr list --state` SHALL be documented and validated. Recommended set: opened, closed, merged. Invalid values SHALL result in exit code 2. So that Agent can discover valid options without trial and error.

#### Scenario: Allowed state values documented
- **WHEN** user or Agent runs `gitlab-tools mr list --help` or consults docs
- **THEN** the allowed values for `--state` SHALL be listed or linked
- **AND** invalid values SHALL result in exit code 2 with a usage or validation error
