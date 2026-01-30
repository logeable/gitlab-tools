## ADDED Requirements

### Requirement: Machine-Readable Output
The system SHALL support optional machine-readable output (e.g. `--json`) for all commands in this capability (branch list, branch diff). When the designated flag is set, the command SHALL output the primary result as JSON to stdout with a stable structure, as specified by the agent-interface capability.

#### Scenario: JSON output for branch commands
- **WHEN** user or Agent runs `gitlab-tools branch list` or `gitlab-tools branch diff <project> <source> <target>` with the machine-output flag (e.g. `--json`)
- **THEN** the system SHALL print valid JSON to stdout with the key fields of the result (e.g. branch name, commit sha, commits, file stats for diff)
- **AND** the output SHALL be parseable by scripts and Agent for chaining or decision-making

### Requirement: Branch List Output for Composition
The system SHALL ensure `branch list [project]` output includes per branch: name, commit sha (or short_sha), and project identification when listing all projects. So that Agent can chain to `branch diff <project> <source> <target>` or `mr create <project> <source> <target>` using branch names from the list.

#### Scenario: List output includes branch name and commit
- **WHEN** user or Agent runs `gitlab-tools branch list <project>` with machine-readable output
- **THEN** each branch SHALL include at least name and commit identifier (sha or short_sha)
- **AND** when no project is given (list all projects), each group SHALL include project id or path so that branch names can be used with that project in diff or mr create

### Requirement: Branch Diff as Atomic Compare
The system SHALL provide `branch diff <project> <source_branch> <target_branch>` as a single atomic operation. It SHALL accept exactly one project identifier and two branch names. Output SHALL include commits and file statistics; flags `--stat` and `--commits` SHALL restrict output to a subset only, without changing the API call. Optional `--create-mr` SHALL be a composition: perform diff then create MR; MR creation SHALL use the same project and branch names.

#### Scenario: Diff output supports decision or MR creation
- **WHEN** user or Agent runs `gitlab-tools branch diff <project> <source> <target>` with machine-readable output
- **THEN** the output SHALL include at least source_branch, target_branch, commit count or list, and file change summary (adds, modifies, deletes)
- **AND** so that Agent can decide to run `mr create <project> <source> <target>` with the same arguments if desired

#### Scenario: Create MR after diff is explicit composition
- **WHEN** user or Agent runs `branch diff ... --create-mr [--mr-title ...] [--mr-description ...]`
- **THEN** the system SHALL first perform the branch comparison, then create one MR from source to target
- **AND** the MR creation SHALL be equivalent to calling `mr create <project> <source> <target>` with the same options; output SHALL include MR iid and web_url for chaining (e.g. to mr merge)

### Requirement: Branch List Without Project (All Projects)
When no project argument is given, `branch list` SHALL list branches for all projects accessible to the user, grouped by project with project identification (id or path) per group. Documentation SHALL note that this can be slow or return a large result and SHALL recommend passing a project when known so that Agent can prefer scoped calls for performance.

#### Scenario: List all projects' branches
- **WHEN** user or Agent runs `gitlab-tools branch list` without a project argument
- **THEN** the system SHALL fetch branches for each accessible project and SHALL display results grouped by project
- **AND** each group SHALL include project id or path so that branch names can be used with that project in diff or mr create

#### Scenario: Docs recommend scoping by project
- **WHEN** user or Agent consults docs or help for `branch list`
- **THEN** the documentation SHALL note that omitting project lists all projects and may be slow or large
- **AND** SHALL recommend passing a project (id or path) when known for faster and smaller output

### Requirement: Branch List Quiet and Hide-Empty Modes
The system SHALL support `branch list` with optional `--quiet` (output only project names, one per line) and `--hide-empty` (hide projects with no branches). The output shape SHALL be documented so that Agent knows what to expect: with `--quiet`, one project name per line; with `--hide-empty`, only projects that have at least one branch (or match search) are shown.

#### Scenario: Quiet output shape
- **WHEN** user or Agent runs `gitlab-tools branch list [project] --quiet`
- **THEN** the system SHALL output only project names, one per line
- **AND** no detailed branch information SHALL be printed so that output is suitable for piping or parsing project list only

#### Scenario: Hide-empty excludes empty projects
- **WHEN** user or Agent runs `gitlab-tools branch list --hide-empty` without a project argument
- **THEN** the system SHALL hide projects that have no branches (or whose branches are all filtered out by --search)
- **AND** only projects with at least one branch SHALL appear in the output

### Requirement: Branch List Search Semantics
The system SHALL support `branch list` with optional `--search <term>`. The search SHALL be documented as case-insensitive substring match on branch name. So that Agent can filter branches by name (e.g. all feature/*) without relying on exact match.

#### Scenario: Search is substring and case-insensitive
- **WHEN** user or Agent runs `gitlab-tools branch list [project] --search "feature"`
- **THEN** the system SHALL return only branches whose name contains the term (case-insensitive)
- **AND** the semantics (substring, case-insensitive) SHALL be documented in help or docs so that Agent can predict results
