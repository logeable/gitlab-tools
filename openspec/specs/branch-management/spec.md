# branch-management Specification

## Purpose
提供 GitLab 项目分支的管理功能，包括列出单个项目或所有项目的分支信息，支持按分支名搜索过滤。
## Requirements
### Requirement: List Project Branches
The system SHALL provide a command to list all branches for one or more GitLab projects. The project identifier is optional; if not provided, branches from all accessible projects SHALL be listed.

#### Scenario: List branches for all projects
- **WHEN** user executes `gitlab-tools branch list` without a project identifier
- **THEN** the system SHALL fetch the list of all accessible projects from GitLab API
- **AND** the system SHALL fetch branches for each project
- **AND** the system SHALL display branches grouped by project with project identification
- **AND** the system SHALL display all branches with details including branch name, protected status, default branch indicator, last commit SHA, message, commit author name, and commit author email

#### Scenario: List branches by numeric project ID
- **WHEN** user executes `gitlab-tools branch list <项目ID>` with a numeric project ID
- **THEN** the system SHALL fetch the branch list from GitLab API for the specified project
- **AND** the system SHALL display all branches with details including branch name, protected status, default branch indicator, last commit SHA, message, commit author name, and commit author email

#### Scenario: List branches by project path
- **WHEN** user executes `gitlab-tools branch list <项目路径>` with a project path (e.g., `my-group/my-project`)
- **THEN** the system SHALL fetch the branch list from GitLab API using the path as identifier
- **AND** the system SHALL display all branches in the same format as numeric ID

#### Scenario: Filter branches by search term
- **WHEN** user executes `gitlab-tools branch list [项目ID] --search <搜索词>` with or without a project identifier
- **THEN** the system SHALL fetch the branch list from GitLab API (for specified project or all projects)
- **AND** the system SHALL filter branches whose names contain the search term (case-insensitive partial match)
- **AND** the system SHALL display only the matching branches

#### Scenario: Hide empty projects
- **WHEN** user executes `gitlab-tools branch list --hide-empty` without a project identifier
- **THEN** the system SHALL fetch branches for all accessible projects
- **AND** the system SHALL hide projects that have no branches (or all branches are filtered out by search)
- **AND** the system SHALL display only projects that have at least one branch

#### Scenario: Quiet mode output
- **WHEN** user executes `gitlab-tools branch list [项目ID] --quiet`
- **THEN** the system SHALL fetch the branch list from GitLab API
- **AND** the system SHALL display only project names (one per line)
- **AND** the system SHALL display project names only for projects that have at least one branch
- **AND** the system SHALL not display detailed branch information

#### Scenario: Combine quiet and hide-empty
- **WHEN** user executes `gitlab-tools branch list --quiet --hide-empty`
- **THEN** the system SHALL fetch branches for all accessible projects
- **AND** the system SHALL display only project names for projects that have branches
- **AND** the system SHALL hide projects with no branches

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle authentication errors
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

#### Scenario: Handle empty branch list
- **WHEN** the project has no branches or all branches are filtered out by search
- **AND** `--hide-empty` is not specified
- **THEN** the system SHALL display a message indicating no branches were found

### Requirement: Compare Branch Differences
The system SHALL provide a command to compare differences between two branches in a GitLab project. The command SHALL display commit differences and file change statistics between the source branch and target branch.

#### Scenario: Compare branches with full diff information
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支>` with a project identifier and two branch names
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display project identifier, source branch name, and target branch name
- **AND** the system SHALL display commit differences including commit SHA, author name, commit message, and commit time
- **AND** the system SHALL display file change statistics including number of files added, modified, and deleted
- **AND** the system SHALL display detailed file changes including file paths and change types (added/modified/deleted/renamed)

#### Scenario: Compare branches by project path
- **WHEN** user executes `gitlab-tools branch diff <项目路径> <源分支> <目标分支>` with a project path (e.g., `my-group/my-project`) and two branch names
- **THEN** the system SHALL fetch the branch comparison from GitLab API using the path as identifier
- **AND** the system SHALL display the comparison results in the same format as numeric project ID

#### Scenario: Display only file statistics
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支> --stat`
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display only file change statistics (number of files added, modified, deleted)
- **AND** the system SHALL not display commit details or detailed file changes

#### Scenario: Display only commit list
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支> --commits`
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display only commit differences (commit SHA, author, message, time)
- **AND** the system SHALL not display file statistics or detailed file changes

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle invalid branch names
- **WHEN** user provides invalid or non-existent branch names
- **THEN** the system SHALL return an error message indicating the branch could not be found

#### Scenario: Handle authentication errors
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

#### Scenario: Compare same branch
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <分支名> <分支名>` with the same branch name for both source and target
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display a message indicating no differences or empty results

#### Scenario: Create merge request after diff
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支> --create-mr`
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display the branch differences
- **AND** the system SHALL create a Merge Request from source branch to target branch
- **AND** the system SHALL display the created Merge Request information including MR ID, title, and Web URL

#### Scenario: Create merge request with custom title
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支> --create-mr --mr-title "自定义标题"`
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display the branch differences
- **AND** the system SHALL create a Merge Request with the specified title
- **AND** the system SHALL display the created Merge Request information

#### Scenario: Create merge request with title and description
- **WHEN** user executes `gitlab-tools branch diff <项目ID> <源分支> <目标分支> --create-mr --mr-title "标题" --mr-description "描述"`
- **THEN** the system SHALL fetch the branch comparison from GitLab API
- **AND** the system SHALL display the branch differences
- **AND** the system SHALL create a Merge Request with the specified title and description
- **AND** the system SHALL display the created Merge Request information

#### Scenario: Handle merge request creation errors
- **WHEN** user attempts to create a Merge Request but the MR already exists or there are permission issues
- **THEN** the system SHALL return an appropriate error message indicating the reason for failure

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

