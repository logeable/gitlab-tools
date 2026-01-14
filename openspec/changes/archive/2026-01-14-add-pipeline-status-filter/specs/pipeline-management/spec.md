# pipeline-management Specification

## Purpose
提供 GitLab 项目 Pipeline 的管理功能，包括查看 Pipeline 详细信息和列出项目的 Pipeline 列表，支持按状态过滤。

## ADDED Requirements
### Requirement: List Project Pipelines
The system SHALL provide a command to list pipelines for one or more GitLab projects.

#### Scenario: List pipelines for single project by numeric ID
- **WHEN** user executes `gitlab-tools pipeline list <项目ID>` with a numeric project ID
- **THEN** the system SHALL fetch the list of pipelines from GitLab API for the specified project
- **AND** the system SHALL display all pipelines with details including pipeline ID, status, ref, SHA, source, created time, updated time, and Web URL

#### Scenario: List pipelines for multiple projects
- **WHEN** user executes `gitlab-tools pipeline list <项目ID1> <项目ID2> ...` with multiple project IDs
- **THEN** the system SHALL fetch the list of pipelines for each specified project
- **AND** the system SHALL display pipelines grouped by project with project identification

#### Scenario: List pipelines by project path
- **WHEN** user executes `gitlab-tools pipeline list <项目路径>` with a project path (e.g., `my-group/my-project`)
- **THEN** the system SHALL fetch the list of pipelines from GitLab API using the path as identifier
- **AND** the system SHALL display pipelines in the same format as numeric ID

#### Scenario: List pipelines with status filter
- **WHEN** user executes `gitlab-tools pipeline list <项目ID> --status <状态>` with a project ID and status value (e.g., success, failed, running)
- **THEN** the system SHALL fetch only pipelines with the specified status from GitLab API
- **AND** the system SHALL display only pipelines matching the status filter

#### Scenario: List pipelines with status filter by project path
- **WHEN** user executes `gitlab-tools pipeline list <项目路径> --status <状态>` with a project path and status value
- **THEN** the system SHALL fetch only pipelines with the specified status using the path as identifier
- **AND** the system SHALL display only pipelines matching the status filter

#### Scenario: List pipelines with status and limit
- **WHEN** user executes `gitlab-tools pipeline list <项目ID> --status <状态> --limit <数量>` with a project ID, status value, and limit
- **THEN** the system SHALL fetch pipelines with the specified status and limit the number of results
- **AND** the system SHALL display the filtered and limited pipelines

#### Scenario: Handle invalid status value
- **WHEN** user provides an invalid status value via `--status` parameter
- **THEN** the system SHALL return an appropriate error message indicating the status value is invalid

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle authentication errors
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

#### Scenario: Handle empty pipeline list
- **WHEN** the project has no pipelines or no pipelines match the status filter
- **THEN** the system SHALL display a message indicating no pipelines were found

