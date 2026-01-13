## ADDED Requirements
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

