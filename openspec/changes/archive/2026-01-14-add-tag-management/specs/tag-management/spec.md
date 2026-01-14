# tag-management Specification

## Purpose
提供 GitLab 项目标签（Tag）的管理功能，包括列出项目标签和在指定分支上创建标签。

## ADDED Requirements
### Requirement: List Tags
The system SHALL provide a command to list all tags for a GitLab project.

#### Scenario: List tags by numeric project ID
- **WHEN** user executes `gitlab-tools tag list <项目ID>` with a numeric project ID
- **THEN** the system SHALL fetch the list of tags from GitLab API for the specified project
- **AND** the system SHALL display all tags with details including tag name, commit SHA, commit message, author, created time, and Web URL

#### Scenario: List tags by project path
- **WHEN** user executes `gitlab-tools tag list <项目路径>` with a project path (e.g., `my-group/my-project`)
- **THEN** the system SHALL fetch the list of tags from GitLab API using the path as identifier
- **AND** the system SHALL display all tags in the same format as numeric ID

#### Scenario: Handle empty tag list
- **WHEN** the project has no tags
- **THEN** the system SHALL display a message indicating no tags were found

#### Scenario: Handle invalid project identifier for list
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle authentication errors for list
- **WHEN** user attempts to access a project without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

### Requirement: Create Tag
The system SHALL provide a command to create a tag in a GitLab project.

#### Scenario: Create tag on default branch (main)
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名>` with a numeric project ID and tag name
- **THEN** the system SHALL create a tag on the default branch (main) at the latest commit
- **AND** the system SHALL display the created tag information including tag name, commit SHA, and Web URL

#### Scenario: Create tag on default branch by project path
- **WHEN** user executes `gitlab-tools tag create <项目路径> <标签名>` with a project path (e.g., `my-group/my-project`) and tag name
- **THEN** the system SHALL create a tag on the default branch (main) using the path as identifier
- **AND** the system SHALL display the created tag information in the same format as numeric ID

#### Scenario: Create tag on specified branch
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --branch <分支名>` with a project ID, tag name, and branch name
- **THEN** the system SHALL create a tag on the specified branch at the latest commit of that branch
- **AND** the system SHALL display the created tag information

#### Scenario: Create tag at specific commit
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --ref <提交SHA>` with a project ID, tag name, and commit SHA
- **THEN** the system SHALL create a tag at the specified commit SHA
- **AND** the system SHALL display the created tag information

#### Scenario: Create tag with message
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --message "标签消息"` with a project ID, tag name, and message
- **THEN** the system SHALL create a tag with the specified message
- **AND** the system SHALL display the created tag information including the message

#### Scenario: Create tag with branch and message
- **WHEN** user executes `gitlab-tools tag create <项目ID> <标签名> --branch <分支名> --message "标签消息"` with a project ID, tag name, branch name, and message
- **THEN** the system SHALL create a tag on the specified branch with the specified message
- **AND** the system SHALL display the created tag information

#### Scenario: Handle invalid project identifier
- **WHEN** user provides an invalid or non-existent project identifier
- **THEN** the system SHALL return an error message indicating the project could not be found or accessed

#### Scenario: Handle invalid branch name
- **WHEN** user provides an invalid or non-existent branch name via `--branch` parameter
- **THEN** the system SHALL return an error message indicating the branch could not be found

#### Scenario: Handle invalid commit SHA
- **WHEN** user provides an invalid or non-existent commit SHA via `--ref` parameter
- **THEN** the system SHALL return an error message indicating the commit could not be found

#### Scenario: Handle duplicate tag name
- **WHEN** user attempts to create a tag with a name that already exists
- **THEN** the system SHALL return an appropriate error message indicating the tag already exists

#### Scenario: Handle authentication errors
- **WHEN** user attempts to create a tag without proper authentication or insufficient permissions
- **THEN** the system SHALL return an appropriate error message

