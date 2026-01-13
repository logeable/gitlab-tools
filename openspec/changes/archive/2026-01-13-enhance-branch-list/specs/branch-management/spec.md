# Branch Management Specification

## MODIFIED Requirements

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

