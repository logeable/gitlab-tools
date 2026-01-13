# merge-request-management Specification

## Purpose
提供 GitLab 项目 Merge Request 的管理功能，包括查询开放的 Merge Request 列表、创建 Merge Request 和合并 Merge Request。
## Requirements
### Requirement: List Open Merge Requests
The system SHALL provide a command to list all open merge requests for a GitLab project.

#### Scenario: List open merge requests by numeric project ID
- **WHEN** user executes `gitlab-tools mr list <项目ID>` with a numeric project ID
- **THEN** the system SHALL fetch the list of open merge requests from GitLab API for the specified project
- **AND** the system SHALL display all open merge requests with details including MR IID, title, source branch, target branch, state, author, created time, and Web URL

#### Scenario: List open merge requests by project path
- **WHEN** user executes `gitlab-tools mr list <项目路径>` with a project path (e.g., `my-group/my-project`)
- **THEN** the system SHALL fetch the list of open merge requests from GitLab API using the path as identifier
- **AND** the system SHALL display all open merge requests in the same format as numeric ID

#### Scenario: Handle empty merge request list
- **WHEN** the project has no open merge requests
- **THEN** the system SHALL display a message indicating no open merge requests were found

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle authentication errors
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

### Requirement: Merge Merge Request
The system SHALL provide a command to merge a merge request in a GitLab project.

#### Scenario: Merge merge request by numeric project ID and MR IID
- **WHEN** user executes `gitlab-tools mr merge <项目ID> <MR IID>` with a numeric project ID and MR IID
- **THEN** the system SHALL fetch the merge request from GitLab API
- **AND** the system SHALL merge the merge request using GitLab API
- **AND** the system SHALL display the merge result including merged status and Web URL

#### Scenario: Merge merge request by project path
- **WHEN** user executes `gitlab-tools mr merge <项目路径> <MR IID>` with a project path (e.g., `my-group/my-project`) and MR IID
- **THEN** the system SHALL fetch the merge request from GitLab API using the path as identifier
- **AND** the system SHALL merge the merge request in the same way as numeric project ID
- **AND** the system SHALL display the merge result

#### Scenario: Merge merge request with source branch deletion
- **WHEN** user executes `gitlab-tools mr merge <项目ID> <MR IID> --delete-source-branch`
- **THEN** the system SHALL merge the merge request from GitLab API
- **AND** the system SHALL delete the source branch after successful merge
- **AND** the system SHALL display the merge result

#### Scenario: Merge merge request with custom merge commit message
- **WHEN** user executes `gitlab-tools mr merge <项目ID> <MR IID> --merge-commit-message "自定义消息"`
- **THEN** the system SHALL merge the merge request from GitLab API
- **AND** the system SHALL use the specified merge commit message
- **AND** the system SHALL display the merge result

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle invalid MR IID
- **WHEN** user provides an invalid or non-existent MR IID
- **THEN** the system SHALL return an error message indicating the merge request could not be found

#### Scenario: Handle already merged merge request
- **WHEN** user attempts to merge a merge request that is already merged
- **THEN** the system SHALL return an appropriate error message indicating the merge request is already merged

#### Scenario: Handle merge conflicts
- **WHEN** user attempts to merge a merge request that has conflicts
- **THEN** the system SHALL return an appropriate error message indicating the merge request has conflicts that need to be resolved

#### Scenario: Handle authentication errors
- **WHEN** user attempts to merge a merge request without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

